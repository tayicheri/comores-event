package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery intercepte les panics et répond 500 JSON.
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				slog.Error("panic récupérée",
					"panic", rec,
					"method", r.Method,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				writeError(w, http.StatusInternalServerError, "Erreur interne du serveur")
			}()
			next.ServeHTTP(w, r)
		})
	}
}
