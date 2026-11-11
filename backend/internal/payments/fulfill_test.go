package payments

import (
	"context"
	"sync"
	"testing"
	"time"

	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

type recordingMailer struct {
	mu    sync.Mutex
	calls []*models.Billet
}

func (m *recordingMailer) SendTicketConfirmation(_ context.Context, b *models.Billet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, b)
	return nil
}

func (m *recordingMailer) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
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

func TestFulfillSucceededSendsEmailOnce(t *testing.T) {
	store := &memBillets{}
	mailer := &recordingMailer{}
	meta := BuildMetadata(models.CreateBilletInput{
		FirstName: "Amina", LastName: "Said", Email: "a@example.com", City: "Moroni",
		Quantity: 1, DayVendredi: true, AmountCents: 2000, StripePaymentIntentID: "pi_ok",
	})

	for i := 0; i < 2; i++ {
		b, err := FulfillSucceeded(context.Background(), store, mailer, "pi_ok", meta)
		if err != nil {
			t.Fatalf("iter %d: %v", i, err)
		}
		if b == nil || b.ID != "b1" {
			t.Fatalf("iter %d unexpected billet: %#v", i, b)
		}
	}
	if mailer.count() != 1 {
		t.Fatalf("want 1 email, got %d", mailer.count())
	}
}
