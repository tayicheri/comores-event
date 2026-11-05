package ratelimit

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient parse REDIS_URL et ouvre un client (sans vérifier la connexion).
func NewRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("REDIS_URL invalide: %w", err)
	}
	return redis.NewClient(opts), nil
}

// RedisLimiter : fenêtre fixe via INCR + EXPIRE. En cas d'erreur Redis,
// bascule sur le limiteur de secours (mémoire) pour ne pas bloquer l'API.
type RedisLimiter struct {
	client   *redis.Client
	fallback Limiter
}

// NewRedisLimiter crée un limiteur Redis avec fallback.
func NewRedisLimiter(client *redis.Client, fallback Limiter) *RedisLimiter {
	return &RedisLimiter{client: client, fallback: fallback}
}

func (l *RedisLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	pipe := l.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	ttlCmd := pipe.PTTL(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil {
		slog.Warn("ratelimit: redis indisponible, fallback mémoire", "error", err)
		return l.fallback.Allow(ctx, key, limit, window)
	}

	count := int(incr.Val())
	ttl := ttlCmd.Val()
	// Première requête de la fenêtre, ou clé sans TTL (crash entre INCR et EXPIRE).
	if count == 1 || ttl < 0 {
		if err := l.client.Expire(ctx, key, window).Err(); err != nil {
			slog.Warn("ratelimit: expire a échoué", "error", err)
		}
		ttl = window
	}
	return buildResult(count, limit, ttl), nil
}
