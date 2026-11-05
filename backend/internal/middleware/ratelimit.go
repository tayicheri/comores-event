package middleware

import (
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"comores-event-api/internal/ratelimit"
)

// TrustedProxies détermine quand X-Forwarded-For / X-Real-IP sont dignes de confiance.
// Valeurs de TRUSTED_PROXIES : "false" (aucun), "true" (tous), ou liste d'IP/CIDR séparés par des virgules.
type TrustedProxies struct {
	all      bool
	prefixes []netip.Prefix
}

// ParseTrustedProxies analyse la valeur de TRUSTED_PROXIES.
func ParseTrustedProxies(spec string) *TrustedProxies {
	tp := &TrustedProxies{}
	spec = strings.TrimSpace(strings.ToLower(spec))
	switch spec {
	case "", "false", "0", "no":
		return tp
	case "true", "1", "yes", "*":
		tp.all = true
		return tp
	}
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if p, err := netip.ParsePrefix(part); err == nil {
			tp.prefixes = append(tp.prefixes, p.Masked())
			continue
		}
		if a, err := netip.ParseAddr(part); err == nil {
			tp.prefixes = append(tp.prefixes, netip.PrefixFrom(a, a.BitLen()))
			continue
		}
		slog.Warn("TRUSTED_PROXIES: entrée ignorée", "value", part)
	}
	return tp
}

func (tp *TrustedProxies) trusted(addr netip.Addr) bool {
	if tp == nil {
		return false
	}
	if tp.all {
		return true
	}
	addr = addr.Unmap()
	for _, p := range tp.prefixes {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

// ClientIP retourne l'IP du client en tenant compte des proxys de confiance.
func (tp *TrustedProxies) ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	remote, err := netip.ParseAddr(host)
	if err != nil || !tp.trusted(remote) {
		return host
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		// On remonte la chaîne depuis la droite jusqu'au premier hop non fiable.
		for i := len(parts) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(parts[i])
			a, err := netip.ParseAddr(ip)
			if err != nil {
				break
			}
			if !tp.trusted(a) || i == 0 {
				return a.Unmap().String()
			}
		}
	}
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		if a, err := netip.ParseAddr(xr); err == nil {
			return a.Unmap().String()
		}
	}
	return host
}

// RateLimit limite à `limit` requêtes par `window` et par IP pour le bucket `name`.
// limit <= 0 désactive la limitation.
func RateLimit(l ratelimit.Limiter, tp *TrustedProxies, name string, limit int, window time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		if limit <= 0 || l == nil {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := "rl:" + name + ":" + tp.ClientIP(r)
			res, err := l.Allow(r.Context(), key, limit, window)
			if err != nil {
				// Fail-open : un incident du limiteur ne doit pas couper l'API.
				slog.Warn("ratelimit: erreur", "error", err)
				next.ServeHTTP(w, r)
				return
			}
			h := w.Header()
			h.Set("X-RateLimit-Limit", strconv.Itoa(limit))
			h.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
			if !res.Allowed {
				retry := int(math.Ceil(res.RetryAfter.Seconds()))
				if retry < 1 {
					retry = 1
				}
				h.Set("Retry-After", strconv.Itoa(retry))
				writeError(w, http.StatusTooManyRequests, "Trop de requêtes, veuillez réessayer plus tard")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
