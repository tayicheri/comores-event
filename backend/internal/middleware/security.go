package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"comores-event-api/internal/constants"
)

// SecurityHeaders ajoute les en-têtes de sécurité HTTP standards.
func SecurityHeaders(isProduction bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
			if isProduction {
				h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			}
			if strings.HasPrefix(r.URL.Path, "/api/") {
				h.Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CSRF applique le double-submit cookie : pour les méthodes non sûres, l'en-tête
// X-XSRF-TOKEN doit égaler le cookie XSRF-TOKEN. login/register/refresh sont exemptés.
func CSRF() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
				next.ServeHTTP(w, r)
				return
			}
			if _, exempt := constants.CSRFExemptPaths[r.URL.Path]; exempt {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(constants.CookieCSRF)
			header := r.Header.Get(constants.HeaderCSRF)
			if err != nil || cookie.Value == "" || header == "" ||
				subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(header)) != 1 {
				writeError(w, http.StatusForbidden, "Jeton CSRF invalide ou manquant")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
