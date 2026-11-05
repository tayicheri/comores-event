package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"comores-event-api/internal/models"
	"comores-event-api/internal/payments"
	"comores-event-api/internal/repository"
)

type fakeStripe struct {
	createID     string
	createSecret string
	createStatus string
	createErr    error

	retrieveStatus string
	retrieveAmount int
	retrieveMeta   map[string]string
	retrieveErr    error
}

func (f fakeStripe) CreatePaymentIntent(amount int, currency, receiptEmail string, metadata map[string]string) (string, string, string, error) {
	if f.createErr != nil {
		return "", "", "", f.createErr
	}
	return f.createID, f.createSecret, f.createStatus, nil
}

func (f fakeStripe) RetrievePaymentIntent(id string) (string, int, string, map[string]string, error) {
	if f.retrieveErr != nil {
		return "", 0, "", nil, f.retrieveErr
	}
	return f.retrieveStatus, f.retrieveAmount, "eur", f.retrieveMeta, nil
}

type memBillets struct {
	byPI map[string]*models.Billet
}

func (m *memBillets) CreateIfAbsent(_ context.Context, in models.CreateBilletInput) (*models.Billet, bool, error) {
	if m.byPI == nil {
		m.byPI = map[string]*models.Billet{}
	}
	if existing, ok := m.byPI[in.StripePaymentIntentID]; ok {
		return existing, false, nil
	}
	b := &models.Billet{
		ID:                    "b1",
		PurchasedAt:           time.Now().UTC(),
		FirstName:             in.FirstName,
		LastName:              in.LastName,
		Email:                 in.Email,
		City:                  in.City,
		Quantity:              in.Quantity,
		DayVendredi:           in.DayVendredi,
		DaySamedi:             in.DaySamedi,
		Status:                "valid",
		AmountCents:           in.AmountCents,
		StripePaymentIntentID: in.StripePaymentIntentID,
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             time.Now().UTC(),
	}
	m.byPI[in.StripePaymentIntentID] = b
	return b, true, nil
}

func (m *memBillets) GetByPaymentIntent(_ context.Context, paymentIntentID string) (*models.Billet, error) {
	b, ok := m.byPI[paymentIntentID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return b, nil
}

func (m *memBillets) List(context.Context) ([]models.Billet, error) { return nil, nil }
func (m *memBillets) Stats(context.Context) (models.BilletStats, error) {
	return models.BilletStats{}, nil
}
func (m *memBillets) Update(context.Context, string, *string, *string) (*models.Billet, error) {
	return nil, repository.ErrNotFound
}

func TestCreateIntent(t *testing.T) {
	h := NewPaymentHandler(
		&memBillets{},
		fakeStripe{createID: "pi_1", createSecret: "sec_1", createStatus: "requires_payment_method"},
		payments.PricingConfig{DayCents: 2000, PassCents: 3000},
		"",
	)
	body := `{"first_name":"Amina","last_name":"Said","email":"a@example.com","city":"Moroni","formula":"day","day":"vendredi","quantity":2}`
	req := httptest.NewRequest(http.MethodPost, "/api/payments/intent", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreateIntent(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["payment_intent_id"] != "pi_1" || resp["amount_cents"].(float64) != 4000 {
		t.Fatalf("unexpected resp: %#v", resp)
	}
}

func TestConfirmIdempotent(t *testing.T) {
	store := &memBillets{}
	meta := payments.BuildMetadata(models.CreateBilletInput{
		FirstName: "Amina", LastName: "Said", Email: "a@example.com", City: "Moroni",
		Quantity: 1, DayVendredi: true, AmountCents: 2000, StripePaymentIntentID: "pi_ok",
	})
	h := NewPaymentHandler(
		store,
		fakeStripe{
			retrieveStatus: "succeeded",
			retrieveAmount: 2000,
			retrieveMeta:   meta,
		},
		payments.PricingConfig{DayCents: 2000, PassCents: 3000},
		"",
	)
	body := `{"payment_intent_id":"pi_ok"}`
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/payments/confirm", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Confirm(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("iter %d status %d body %s", i, rec.Code, rec.Body.String())
		}
	}
	if len(store.byPI) != 1 {
		t.Fatalf("want 1 billet, got %d", len(store.byPI))
	}
}
