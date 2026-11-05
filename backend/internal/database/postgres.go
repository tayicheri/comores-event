// Package database ouvre la connexion PostgreSQL (sqlx + pgx/v5 stdlib).
package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"comores-event-api/internal/config"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // driver "pgx"
	"github.com/jmoiron/sqlx"
)

// BuildDSN construit une URL de connexion PostgreSQL correctement échappée.
func BuildDSN(c config.DBConfig) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Pass),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     "/" + c.Name,
		RawQuery: url.Values{"sslmode": []string{c.SSLMode}}.Encode(),
	}
	return u.String()
}

// Connect ouvre le pool et vérifie la connexion (avec quelques tentatives).
func Connect(ctx context.Context, c config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", BuildDSN(c))
	if err != nil {
		return nil, fmt.Errorf("ouverture postgres: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	var lastErr error
	for attempt := 1; attempt <= 10; attempt++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		lastErr = db.PingContext(pingCtx)
		cancel()
		if lastErr == nil {
			return db, nil
		}
		select {
		case <-ctx.Done():
			_ = db.Close()
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
		}
	}
	_ = db.Close()
	return nil, fmt.Errorf("ping postgres: %w", lastErr)
}

// IsUniqueViolation indique une violation de contrainte UNIQUE (23505).
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
