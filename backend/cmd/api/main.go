package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "time/tzdata" // fuseaux horaires embarqués (images alpine minimales)

	"comores-event-api/internal/api"
	"comores-event-api/internal/auth"
	"comores-event-api/internal/config"
	"comores-event-api/internal/database"
	"comores-event-api/internal/logging"
	"comores-event-api/internal/middleware"
	"comores-event-api/internal/ratelimit"
	"comores-event-api/internal/repository"

	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("arrêt sur erreur", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		logging.Init("development")
		return err
	}
	logging.Init(cfg.AppEnv)

	if loc, err := time.LoadLocation(cfg.TZ); err == nil {
		time.Local = loc
	} else {
		slog.Warn("TZ invalide, UTC conservé", "tz", cfg.TZ, "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()
	slog.Info("postgres connecté", "host", cfg.DB.Host, "db", cfg.DB.Name)

	// Rate limiting : Redis si REDIS_URL, sinon mémoire.
	memLimiter := ratelimit.NewMemoryLimiter(ctx)
	var (
		limiter ratelimit.Limiter = memLimiter
		rdb     *redis.Client
	)
	if cfg.RedisURL != "" {
		rdb, err = ratelimit.NewRedisClient(cfg.RedisURL)
		if err != nil {
			return err
		}
		defer rdb.Close()
		limiter = ratelimit.NewRedisLimiter(rdb, memLimiter)
		slog.Info("rate limiting via Redis")
	} else {
		slog.Warn("REDIS_URL vide : rate limiting en mémoire (instance unique)")
	}

	users := repository.NewUserRepository(db)
	sessions := repository.NewSessionRepository(db)
	billets := repository.NewBilletRepository(db)

	go purgeExpiredTokens(ctx, sessions)

	handler := api.NewRouter(api.Deps{
		Config:   cfg,
		DB:       db,
		Redis:    rdb,
		Limiter:  limiter,
		Users:    users,
		Sessions: sessions,
		Billets:  billets,
		JWT:      auth.NewJWTManager(cfg.JWTSecret),
		Proxies:  middleware.ParseTrustedProxies(cfg.TrustedProxy),
	})

	srv := &http.Server{
		Addr:              net.JoinHostPort("", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("API démarrée", "addr", srv.Addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("arrêt en cours…")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// purgeExpiredTokens supprime périodiquement les refresh tokens périmés.
func purgeExpiredTokens(ctx context.Context, sessions repository.SessionRepository) {
	t := time.NewTicker(6 * time.Hour)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if n, err := sessions.DeleteExpired(ctx); err != nil {
				slog.Warn("purge refresh_tokens", "error", err)
			} else if n > 0 {
				slog.Info("refresh_tokens purgés", "count", n)
			}
		}
	}
}
