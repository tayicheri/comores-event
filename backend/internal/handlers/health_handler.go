package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// HealthHandler expose /health (liveness) et /health/ready (DB + Redis).
type HealthHandler struct {
	db    *sqlx.DB
	redis *redis.Client // nil si REDIS_URL est vide
}

// NewHealthHandler construit le handler de santé.
func NewHealthHandler(db *sqlx.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: rdb}
}

// Live : GET /health
func (h *HealthHandler) Live(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready : GET /health/ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{}
	ok := true

	if err := h.db.PingContext(ctx); err != nil {
		checks["database"] = "down"
		ok = false
	} else {
		checks["database"] = "ok"
	}

	switch {
	case h.redis == nil:
		checks["redis"] = "disabled"
	default:
		if err := h.redis.Ping(ctx).Err(); err != nil {
			checks["redis"] = "down"
			ok = false
		} else {
			checks["redis"] = "ok"
		}
	}

	status, label := http.StatusOK, "ready"
	if !ok {
		status, label = http.StatusServiceUnavailable, "unavailable"
	}
	WriteJSON(w, status, map[string]any{"status": label, "checks": checks})
}
