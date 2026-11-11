package payments

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"comores-event-api/internal/mail"
	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

const (
	metaFirstName    = "first_name"
	metaLastName     = "last_name"
	metaEmail        = "email"
	metaCity         = "city"
	metaQuantity     = "quantity"
	metaDayVendredi  = "day_vendredi"
	metaDaySamedi    = "day_samedi"
	metaAmountCents  = "amount_cents"
)

// BuildMetadata encode les infos acheteur dans le PaymentIntent.
func BuildMetadata(in models.CreateBilletInput) map[string]string {
	return map[string]string{
		metaFirstName:   in.FirstName,
		metaLastName:    in.LastName,
		metaEmail:       in.Email,
		metaCity:        in.City,
		metaQuantity:    strconv.Itoa(in.Quantity),
		metaDayVendredi: strconv.FormatBool(in.DayVendredi),
		metaDaySamedi:   strconv.FormatBool(in.DaySamedi),
		metaAmountCents: strconv.Itoa(in.AmountCents),
	}
}

// InputFromMetadata reconstruit CreateBilletInput depuis les metadata Stripe.
func InputFromMetadata(paymentIntentID string, meta map[string]string) (models.CreateBilletInput, error) {
	qty, err := strconv.Atoi(strings.TrimSpace(meta[metaQuantity]))
	if err != nil || qty < 1 {
		return models.CreateBilletInput{}, fmt.Errorf("metadata quantity invalide")
	}
	amount, err := strconv.Atoi(strings.TrimSpace(meta[metaAmountCents]))
	if err != nil || amount < 0 {
		return models.CreateBilletInput{}, fmt.Errorf("metadata amount_cents invalide")
	}
	first := strings.TrimSpace(meta[metaFirstName])
	last := strings.TrimSpace(meta[metaLastName])
	email := strings.TrimSpace(meta[metaEmail])
	city := strings.TrimSpace(meta[metaCity])
	if first == "" || last == "" || email == "" || city == "" {
		return models.CreateBilletInput{}, fmt.Errorf("metadata acheteur incomplètes")
	}
	vendredi := meta[metaDayVendredi] == "true"
	samedi := meta[metaDaySamedi] == "true"
	if !vendredi && !samedi {
		return models.CreateBilletInput{}, fmt.Errorf("metadata jours invalides")
	}
	return models.CreateBilletInput{
		FirstName:             first,
		LastName:              last,
		Email:                 email,
		City:                  city,
		Quantity:              qty,
		DayVendredi:           vendredi,
		DaySamedi:             samedi,
		AmountCents:           amount,
		StripePaymentIntentID: paymentIntentID,
	}, nil
}

// FulfillSucceeded crée le billet si le PI est succeeded (idempotent).
// Envoie un email de confirmation uniquement lors d'une création effective.
func FulfillSucceeded(
	ctx context.Context,
	billets repository.BilletRepository,
	mailer mail.Mailer,
	paymentIntentID string,
	meta map[string]string,
) (*models.Billet, error) {
	in, err := InputFromMetadata(paymentIntentID, meta)
	if err != nil {
		return nil, err
	}
	b, created, err := billets.CreateIfAbsent(ctx, in)
	if err != nil {
		return nil, err
	}
	if created && mailer != nil {
		if sendErr := mailer.SendTicketConfirmation(ctx, b); sendErr != nil {
			slog.Error("email confirmation billet", "error", sendErr, "billet_id", b.ID, "pi", paymentIntentID)
		}
	}
	return b, nil
}
