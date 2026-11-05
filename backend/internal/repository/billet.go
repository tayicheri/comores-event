package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/database"
	"comores-event-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// BilletRepository accès PostgreSQL aux billets.
type BilletRepository interface {
	CreateIfAbsent(ctx context.Context, in models.CreateBilletInput) (*models.Billet, bool, error)
	GetByPaymentIntent(ctx context.Context, paymentIntentID string) (*models.Billet, error)
	List(ctx context.Context) ([]models.Billet, error)
	Stats(ctx context.Context) (models.BilletStats, error)
	Update(ctx context.Context, id string, status *string, comments *string) (*models.Billet, error)
}

type billetRepo struct {
	db *sqlx.DB
}

// NewBilletRepository retourne l'implémentation PostgreSQL.
func NewBilletRepository(db *sqlx.DB) BilletRepository {
	return &billetRepo{db: db}
}

const billetColumns = `id, purchased_at, first_name, last_name, email, city, quantity,
	day_vendredi, day_samedi, comments, status::text AS status, amount_cents,
	stripe_payment_intent_id, created_at, updated_at`

// CreateIfAbsent insère un billet ; si le PI existe déjà, retourne l'existant (created=false).
func (r *billetRepo) CreateIfAbsent(ctx context.Context, in models.CreateBilletInput) (*models.Billet, bool, error) {
	var b models.Billet
	err := r.db.GetContext(ctx, &b,
		`INSERT INTO billets (
			first_name, last_name, email, city, quantity,
			day_vendredi, day_samedi, status, amount_cents, stripe_payment_intent_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING `+billetColumns,
		in.FirstName, in.LastName, in.Email, in.City, in.Quantity,
		in.DayVendredi, in.DaySamedi, constants.BilletStatusValid, in.AmountCents, in.StripePaymentIntentID,
	)
	if err == nil {
		return &b, true, nil
	}
	if database.IsUniqueViolation(err) {
		existing, getErr := r.GetByPaymentIntent(ctx, in.StripePaymentIntentID)
		if getErr != nil {
			return nil, false, getErr
		}
		return existing, false, nil
	}
	return nil, false, fmt.Errorf("insert billet: %w", err)
}

func (r *billetRepo) GetByPaymentIntent(ctx context.Context, paymentIntentID string) (*models.Billet, error) {
	var b models.Billet
	err := r.db.GetContext(ctx, &b,
		`SELECT `+billetColumns+` FROM billets WHERE stripe_payment_intent_id = $1`, paymentIntentID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *billetRepo) List(ctx context.Context) ([]models.Billet, error) {
	var rows []models.Billet
	err := r.db.SelectContext(ctx, &rows,
		`SELECT `+billetColumns+` FROM billets ORDER BY purchased_at DESC`)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []models.Billet{}
	}
	return rows, nil
}

func (r *billetRepo) Stats(ctx context.Context) (models.BilletStats, error) {
	var s models.BilletStats
	err := r.db.GetContext(ctx, &s, `
		SELECT
			COUNT(*)::int AS total_count,
			COUNT(*) FILTER (WHERE status = 'valid')::int AS valid_count,
			COUNT(*) FILTER (WHERE status = 'used')::int AS used_count,
			COALESCE(SUM(amount_cents), 0)::int AS revenue_cents,
			COALESCE(SUM(quantity) FILTER (WHERE day_vendredi), 0)::int AS vendredi_places,
			COALESCE(SUM(quantity) FILTER (WHERE day_samedi), 0)::int AS samedi_places
		FROM billets`)
	return s, err
}

func (r *billetRepo) Update(ctx context.Context, id string, status *string, comments *string) (*models.Billet, error) {
	if status == nil && comments == nil {
		return r.getByID(ctx, id)
	}
	var b models.Billet
	err := r.db.GetContext(ctx, &b, `
		UPDATE billets SET
			status = CASE WHEN $2::text IS NOT NULL THEN $2::billet_status ELSE status END,
			comments = CASE WHEN $3::boolean THEN $4 ELSE comments END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING `+billetColumns,
		id,
		status,
		comments != nil,
		comments,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *billetRepo) getByID(ctx context.Context, id string) (*models.Billet, error) {
	var b models.Billet
	err := r.db.GetContext(ctx, &b, `SELECT `+billetColumns+` FROM billets WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}
