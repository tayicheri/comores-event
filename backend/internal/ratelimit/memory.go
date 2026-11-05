// Package ratelimit : limiteurs de débit (Redis INCR+EXPIRE, fallback mémoire).
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Result est le résultat d'une vérification de limite.
type Result struct {
	Allowed    bool
	Remaining  int
	RetryAfter time.Duration
}

// Limiter vérifie et consomme un jeton pour une clé donnée (fenêtre fixe).
type Limiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (Result, error)
}

type memEntry struct {
	count   int
	resetAt time.Time
}

// MemoryLimiter : limiteur en mémoire (processus unique).
type MemoryLimiter struct {
	mu      sync.Mutex
	entries map[string]*memEntry
}

// NewMemoryLimiter crée un limiteur mémoire ; le nettoyage s'arrête avec ctx.
func NewMemoryLimiter(ctx context.Context) *MemoryLimiter {
	m := &MemoryLimiter{entries: make(map[string]*memEntry)}
	go m.janitor(ctx)
	return m
}

func (m *MemoryLimiter) Allow(_ context.Context, key string, limit int, window time.Duration) (Result, error) {
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()

	e, ok := m.entries[key]
	if !ok || !e.resetAt.After(now) {
		e = &memEntry{resetAt: now.Add(window)}
		m.entries[key] = e
	}
	e.count++
	return buildResult(e.count, limit, time.Until(e.resetAt)), nil
}

func (m *MemoryLimiter) janitor(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			m.mu.Lock()
			for k, e := range m.entries {
				if !e.resetAt.After(now) {
					delete(m.entries, k)
				}
			}
			m.mu.Unlock()
		}
	}
}

func buildResult(count, limit int, ttl time.Duration) Result {
	if ttl < 0 {
		ttl = 0
	}
	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}
	if count > limit {
		return Result{Allowed: false, Remaining: 0, RetryAfter: ttl}
	}
	return Result{Allowed: true, Remaining: remaining, RetryAfter: ttl}
}
