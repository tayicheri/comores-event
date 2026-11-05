// Package repository contient les accès PostgreSQL.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/database"
	"comores-event-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// Erreurs métier du repository.
var (
	ErrNotFound     = errors.New("ressource introuvable")
	ErrEmailTaken   = errors.New("email déjà utilisé")
	ErrTokenInvalid = errors.New("refresh token invalide ou expiré")
	ErrTokenReused  = errors.New("refresh token déjà utilisé")
)

// UserRepository décrit l'accès aux utilisateurs.
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string, displayName *string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	UpdateDisplayName(ctx context.Context, id string, displayName *string) error
	BumpTokenVersion(ctx context.Context, id string) error
	RecordFailedLogin(ctx context.Context, id string) error
	ResetFailedLogins(ctx context.Context, id string) error
}

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepository retourne l'implémentation PostgreSQL.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

const userColumns = `id, email, password_hash, display_name, account_status::text AS account_status,
	email_verified_at, token_version, failed_login_count, locked_until,
	created_at, updated_at, deleted_at`

func (r *userRepo) Create(ctx context.Context, email, passwordHash string, displayName *string) (*models.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var u models.User
	err = tx.GetContext(ctx, &u,
		`INSERT INTO users (email, password_hash, display_name)
		 VALUES ($1, $2, $3)
		 RETURNING `+userColumns, email, passwordHash, displayName)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO user_roles (user_id, role_id)
		 SELECT $1, id FROM roles WHERE name = $2`, u.ID, constants.RoleUser)
	if err != nil {
		return nil, fmt.Errorf("insert user_role: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, fmt.Errorf("rôle %q introuvable (schéma non initialisé ?)", constants.RoleUser)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	u.Roles = []string{constants.RoleUser}
	return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	return r.finish(ctx, &u, err)
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return r.finish(ctx, &u, err)
}

func (r *userRepo) finish(ctx context.Context, u *models.User, err error) (*models.User, error) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	roles, err := r.loadRoles(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	return u, nil
}

func (r *userRepo) loadRoles(ctx context.Context, userID string) ([]string, error) {
	roles := []string{}
	err := r.db.SelectContext(ctx, &roles,
		`SELECT r.name FROM roles r
		 JOIN user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id = $1
		 ORDER BY r.name`, userID)
	if err != nil {
		return nil, fmt.Errorf("load roles: %w", err)
	}
	return roles, nil
}

func (r *userRepo) UpdateDisplayName(ctx context.Context, id string, displayName *string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET display_name = $2, updated_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`, id, displayName)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepo) BumpTokenVersion(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET token_version = token_version + 1, updated_at = NOW() WHERE id = $1`, id)
	return err
}

// RecordFailedLogin incrémente le compteur d'échecs et verrouille le compte
// temporairement au-delà du seuil.
func (r *userRepo) RecordFailedLogin(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET
		   failed_login_count = CASE WHEN failed_login_count + 1 >= $2 THEN 0 ELSE failed_login_count + 1 END,
		   locked_until = CASE WHEN failed_login_count + 1 >= $2
		                       THEN NOW() + make_interval(secs => $3::double precision)
		                       ELSE locked_until END,
		   updated_at = NOW()
		 WHERE id = $1`,
		id, constants.MaxFailedLogins, (time.Duration(constants.LockDuration)).Seconds())
	return err
}

func (r *userRepo) ResetFailedLogins(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET failed_login_count = 0, locked_until = NULL
		 WHERE id = $1 AND (failed_login_count <> 0 OR locked_until IS NOT NULL)`, id)
	return err
}
