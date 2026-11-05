package models

import "time"

// Billet correspond à la table billets.
type Billet struct {
	ID                     string    `db:"id"`
	PurchasedAt            time.Time `db:"purchased_at"`
	FirstName              string    `db:"first_name"`
	LastName               string    `db:"last_name"`
	Email                  string    `db:"email"`
	City                   string    `db:"city"`
	Quantity               int       `db:"quantity"`
	DayVendredi            bool      `db:"day_vendredi"`
	DaySamedi              bool      `db:"day_samedi"`
	Comments               *string   `db:"comments"`
	Status                 string    `db:"status"`
	AmountCents            int       `db:"amount_cents"`
	StripePaymentIntentID  string    `db:"stripe_payment_intent_id"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
}

// BilletResponse est l'objet billet exposé par l'API.
type BilletResponse struct {
	ID                    string    `json:"id"`
	PurchasedAt           time.Time `json:"purchased_at"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	Email                 string    `json:"email"`
	City                  string    `json:"city"`
	Quantity              int       `json:"quantity"`
	DayVendredi           bool      `json:"day_vendredi"`
	DaySamedi             bool      `json:"day_samedi"`
	Comments              *string   `json:"comments"`
	Status                string    `json:"status"`
	AmountCents           int       `json:"amount_cents"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
}

// ToResponse convertit un Billet en réponse API.
func (b *Billet) ToResponse() BilletResponse {
	return BilletResponse{
		ID:                    b.ID,
		PurchasedAt:           b.PurchasedAt,
		FirstName:             b.FirstName,
		LastName:              b.LastName,
		Email:                 b.Email,
		City:                  b.City,
		Quantity:              b.Quantity,
		DayVendredi:           b.DayVendredi,
		DaySamedi:             b.DaySamedi,
		Comments:              b.Comments,
		Status:                b.Status,
		AmountCents:           b.AmountCents,
		StripePaymentIntentID: b.StripePaymentIntentID,
	}
}

// CreateBilletInput données nécessaires à l'insertion après paiement.
type CreateBilletInput struct {
	FirstName             string
	LastName              string
	Email                 string
	City                  string
	Quantity              int
	DayVendredi           bool
	DaySamedi             bool
	AmountCents           int
	StripePaymentIntentID string
}

// PaymentIntentRequest : corps de POST /api/payments/intent.
type PaymentIntentRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Email     string `json:"email" validate:"required,email,max=254"`
	City      string `json:"city" validate:"required,min=1,max=120"`
	Formula   string `json:"formula" validate:"required,oneof=day pass"`
	Day       string `json:"day" validate:"omitempty,oneof=vendredi samedi"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=20"`
}

// PaymentConfirmRequest : corps de POST /api/payments/confirm.
type PaymentConfirmRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required,min=1,max=255"`
}

// UpdateBilletRequest : corps de PATCH /api/admin/billets/{id}.
type UpdateBilletRequest struct {
	Status   *string `json:"status" validate:"omitempty,oneof=valid used"`
	Comments *string `json:"comments" validate:"omitempty,max=2000"`
}

// BilletStats agrégats pour le dashboard admin.
type BilletStats struct {
	TotalCount     int `db:"total_count" json:"total_count"`
	ValidCount     int `db:"valid_count" json:"valid_count"`
	UsedCount      int `db:"used_count" json:"used_count"`
	RevenueCents   int `db:"revenue_cents" json:"revenue_cents"`
	VendrediPlaces int `db:"vendredi_places" json:"vendredi_places"`
	SamediPlaces   int `db:"samedi_places" json:"samedi_places"`
}

// BilletListResponse liste + stats.
type BilletListResponse struct {
	Billets []BilletResponse `json:"billets"`
	Stats   BilletStats      `json:"stats"`
}
