// Package api assemble le routeur HTTP et la chaîne de middlewares.
package api

import (
	"net/http"
	"time"

	"comores-event-api/internal/auth"
	"comores-event-api/internal/config"
	"comores-event-api/internal/handlers"
	"comores-event-api/internal/middleware"
	"comores-event-api/internal/ratelimit"
	"comores-event-api/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

// Deps regroupe les dépendances nécessaires au routeur.
type Deps struct {
	Config   *config.Config
	DB       *sqlx.DB
	Redis    *redis.Client // peut être nil
	Limiter  ratelimit.Limiter
	Users    repository.UserRepository
	Sessions repository.SessionRepository
	JWT      *auth.JWTManager
	Proxies  *middleware.TrustedProxies
}

// rateWindow : les limites RATE_LIMIT_* sont exprimées en requêtes par minute et par IP.
const rateWindow = time.Minute

// NewRouter construit le handler racine.
// Ordre des middlewares (externe → interne) :
// Recovery → Logger → Metrics → Cors → CSRF → SecurityHeaders → mux
func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandler(d.DB, d.Redis)
	mux.HandleFunc("GET /health", health.Live)
	mux.HandleFunc("GET /health/ready", health.Ready)
	mux.Handle("GET /metrics", promhttp.Handler())

	registerAuthRoutes(mux, d)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		handlers.WriteError(w, http.StatusNotFound, "Ressource introuvable")
	})

	resolve := func(r *http.Request) string {
		_, pattern := mux.Handler(r)
		if pattern == "/" {
			return "unmatched"
		}
		return pattern
	}

	return middleware.Chain(mux,
		middleware.Recovery(),
		middleware.Logger(d.Proxies),
		middleware.Metrics(resolve),
		middleware.Cors(d.Config.AllowedOrigin),
		middleware.CSRF(),
		middleware.SecurityHeaders(d.Config.IsProduction),
	)
}
