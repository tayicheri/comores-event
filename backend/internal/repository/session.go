package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// SessionRepository décrit l'accès aux refresh tokens (sessions).
type SessionRepository interface {
	// Create insère un nouveau refresh token. Si SessionID est vide, une nouvelle session est créée.
	Create(ctx context.Context, t *models.RefreshToken) error
	// Rotate révoque le token identifié par oldHash et insère next dans la même session.
	Rotate(ctx context.Context, oldHash string, next *models.RefreshToken) error
	// RevokeByHash révoque un refresh token (logout courant).
	RevokeByHash(ctx context.Context, hash string) error
	// RevokeAllForUser révoque tous les refresh tokens d'un utilisateur.
	RevokeAllForUser(ctx context.Context, userID string) error
	// RevokeToken révoque un token précis par id (utilisé en cas d'échec post-rotation).
	RevokeToken(ctx context.Context, id string) error
	// EnforceSessionLimit révoque les sessions les plus anciennes pour garder au plus keep sessions actives.
	EnforceSessionLimit(ctx context.Context, userID string, keep int) error
	// DeleteExpired purge les tokens expirés ou révoqués depuis longtemps.
	DeleteExpired(ctx context.Context) (int64, error)
}

type sessionRepo struct {
	db *sqlx.DB
}

// NewSessionRepository retourne l'implémentation PostgreSQL.
func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepo{db: db}
}

const tokenColumns = `id, user_id, session_id, token_hash, device_name, device_os,
	ip_address, user_agent, last_seen_at, expires_at, revoked_at, created_at`

func insertToken(ctx context.Context, q sqlx.QueryerContext, t *models.RefreshToken) error {
	return q.QueryRowxContext(ctx,
		`INSERT INTO refresh_tokens
		   (user_id, session_id, token_hash, device_name, device_os, ip_address, user_agent, expires_at)
		 VALUES ($1, COALESCE(NULLIF($2::text, '')::uuid, gen_random_uuid()), $3, $4, $5, $6, $7, $8)
		 RETURNING id, session_id, last_seen_at, created_at`,
		t.UserID, t.SessionID, t.TokenHash, t.DeviceName, t.DeviceOS, t.IPAddress, t.UserAgent, t.ExpiresAt,
	).Scan(&t.ID, &t.SessionID, &t.LastSeenAt, &t.CreatedAt)
}

func (r *sessionRepo) Create(ctx context.Context, t *models.RefreshToken) error {
	if err := insertToken(ctx, r.db, t); err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

func (r *sessionRepo) Rotate(ctx context.Context, oldHash string, next *models.RefreshToken) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var old models.RefreshToken
	err = tx.GetContext(ctx, &old,
		`SELECT `+tokenColumns+` FROM refresh_tokens WHERE token_hash = $1 FOR UPDATE`, oldHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenInvalid
		}
		return err
	}

	now := time.Now()
	if old.RevokedAt != nil {
		// Fenêtre de grâce : deux refresh concurrents (onglets multiples).
		if now.Sub(*old.RevokedAt) <= constants.RefreshReuseGrace {
			return ErrTokenInvalid
		}
		// Réutilisation d'un token déjà roté : on révoque toute la session.
		if _, err := tx.ExecContext(ctx,
			`UPDATE refresh_tokens SET revoked_at = NOW()
			 WHERE session_id = $1 AND revoked_at IS NULL`, old.SessionID); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return ErrTokenReused
	}
	if !old.ExpiresAt.After(now) {
		return ErrTokenInvalid
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW(), last_seen_at = NOW() WHERE id = $1`, old.ID); err != nil {
		return err
	}

	next.UserID = old.UserID
	next.SessionID = old.SessionID
	if next.DeviceName == nil {
		next.DeviceName = old.DeviceName
	}
	if next.DeviceOS == nil {
		next.DeviceOS = old.DeviceOS
	}
	if err := insertToken(ctx, tx, next); err != nil {
		return fmt.Errorf("insert rotated token: %w", err)
	}
	return tx.Commit()
}

func (r *sessionRepo) RevokeByHash(ctx context.Context, hash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`, hash)
	return err
}

func (r *sessionRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}

func (r *sessionRepo) RevokeToken(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`, id)
	return err
}

func (r *sessionRepo) EnforceSessionLimit(ctx context.Context, userID string, keep int) error {
	if keep < 0 {
		keep = 0
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW()
		 WHERE user_id = $1 AND revoked_at IS NULL AND session_id IN (
		   SELECT session_id FROM refresh_tokens
		   WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		   GROUP BY session_id
		   ORDER BY MAX(last_seen_at) DESC
		   OFFSET $2
		 )`, userID, keep)
	return err
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens
		 WHERE expires_at < NOW() - INTERVAL '1 day'
		    OR revoked_at < NOW() - INTERVAL '1 day'`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
