package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger journalise chaque requête (méthode, chemin, statut, durée, IP).
func Logger(tp *TrustedProxies) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := newStatusRecorder(w)
			next.ServeHTTP(rec, r)

			level := slog.LevelInfo
			switch {
			case rec.status >= 500:
				level = slog.LevelError
			case rec.status >= 400:
				level = slog.LevelWarn
			}
			// Les sondes de santé / métriques sont bruyantes : niveau debug.
			if r.URL.Path == "/health" || r.URL.Path == "/health/ready" || r.URL.Path == "/metrics" {
				if rec.status < 400 {
					level = slog.LevelDebug
				}
			}
			slog.Log(r.Context(), level, "http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"bytes", rec.bytes,
				"duration_ms", time.Since(start).Milliseconds(),
				"ip", tp.ClientIP(r),
			)
		})
	}
}
