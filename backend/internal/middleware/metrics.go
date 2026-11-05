package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Nombre total de requêtes HTTP.",
	}, []string{"method", "route", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Durée des requêtes HTTP en secondes.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Nombre de requêtes HTTP en cours.",
	})
)

// RouteResolver retourne le pattern de route pour une requête (cardinalité bornée).
type RouteResolver func(*http.Request) string

// Metrics collecte les métriques Prometheus. Le label "route" est le pattern
// du ServeMux (ex. "POST /api/auth/login"), jamais le chemin brut.
func Metrics(resolve RouteResolver) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			httpRequestsInFlight.Inc()
			defer httpRequestsInFlight.Dec()

			rec := newStatusRecorder(w)
			next.ServeHTTP(rec, r)

			route := "unmatched"
			if resolve != nil {
				if p := resolve(r); p != "" {
					route = p
				}
			}
			httpRequestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
			httpRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
		})
	}
}
