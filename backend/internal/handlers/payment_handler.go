package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/mail"
	"comores-event-api/internal/models"
	"comores-event-api/internal/payments"
	"comores-event-api/internal/repository"

	"github.com/stripe/stripe-go/v82"
)

// PaymentHandler gère intent / confirm / webhook Stripe.
type PaymentHandler struct {
	billets       repository.BilletRepository
	stripe        payments.StripeAPI
	pricing       payments.PricingConfig
	mailer        mail.Mailer
	webhookSecret string
}

// NewPaymentHandler construit le handler paiements.
func NewPaymentHandler(
	billets repository.BilletRepository,
	stripe payments.StripeAPI,
	pricing payments.PricingConfig,
	mailer mail.Mailer,
	webhookSecret string,
) *PaymentHandler {
	if mailer == nil {
		mailer = mail.NopMailer{}
	}
	return &PaymentHandler{
		billets:       billets,
		stripe:        stripe,
		pricing:       pricing,
		mailer:        mailer,
		webhookSecret: strings.TrimSpace(webhookSecret),
	}
}

// CreateIntent : POST /api/payments/intent
func (h *PaymentHandler) CreateIntent(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentIntentRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.City = strings.TrimSpace(req.City)

	vendredi, samedi, err := payments.DayFlags(req.Formula, req.Day)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	amount, err := payments.AmountCents(h.pricing, req.Formula, req.Quantity)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	in := models.CreateBilletInput{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		City:         req.City,
		Quantity:     req.Quantity,
		DayVendredi:  vendredi,
		DaySamedi:    samedi,
		AmountCents:  amount,
	}
	meta := payments.BuildMetadata(in)

	id, secret, status, err := h.stripe.CreatePaymentIntent(amount, constants.CurrencyEUR, req.Email, meta)
	if err != nil {
		slog.Error("stripe: create payment intent", "error", err)
		WriteError(w, http.StatusBadGateway, "Impossible de démarrer le paiement")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"payment_intent_id": id,
		"client_secret":     secret,
		"status":            status,
		"amount_cents":      amount,
		"currency":          constants.CurrencyEUR,
	})
}

// Confirm : POST /api/payments/confirm
func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentConfirmRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	piID := strings.TrimSpace(req.PaymentIntentID)

	status, amount, _, meta, err := h.stripe.RetrievePaymentIntent(piID)
	if err != nil {
		slog.Error("stripe: retrieve payment intent", "error", err, "pi", piID)
		WriteError(w, http.StatusBadGateway, "Impossible de vérifier le paiement")
		return
	}
	if status != string(stripe.PaymentIntentStatusSucceeded) {
		WriteError(w, http.StatusConflict, "Le paiement n'est pas encore validé")
		return
	}

	// Cohérence montant metadata vs PI Stripe.
	if raw, ok := meta["amount_cents"]; ok && raw != "" {
		in, metaErr := payments.InputFromMetadata(piID, meta)
		if metaErr == nil && in.AmountCents != amount {
			slog.Error("stripe: montant incohérent", "meta", in.AmountCents, "pi", amount)
			WriteError(w, http.StatusBadGateway, "Montant de paiement incohérent")
			return
		}
	}

	billet, err := payments.FulfillSucceeded(r.Context(), h.billets, h.mailer, piID, meta)
	if err != nil {
		slog.Error("fulfill billet", "error", err, "pi", piID)
		WriteError(w, http.StatusInternalServerError, "Paiement reçu mais billet non créé")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"billet": billet.ToResponse()})
}

// StripeWebhook : POST /api/webhooks/stripe
func (h *PaymentHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	if h.webhookSecret == "" {
		WriteError(w, http.StatusServiceUnavailable, "Webhook Stripe non configuré")
		return
	}
	const maxBody = 1 << 20
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Corps invalide")
		return
	}
	event, err := payments.ConstructStripeEvent(payload, r.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Signature Stripe invalide")
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			WriteError(w, http.StatusBadRequest, "Événement invalide")
			return
		}
		meta := map[string]string{}
		for k, v := range pi.Metadata {
			meta[k] = v
		}
		if _, err := payments.FulfillSucceeded(r.Context(), h.billets, h.mailer, pi.ID, meta); err != nil {
			slog.Error("webhook fulfill", "error", err, "pi", pi.ID)
			WriteError(w, http.StatusInternalServerError, "Échec création billet")
			return
		}
	default:
		// Ignorer les autres événements.
	}
	WriteJSON(w, http.StatusOK, map[string]string{"received": "ok"})
}

// ErrStripeNotConfigured signal config manquante (tests / health optionnel).
var ErrStripeNotConfigured = errors.New("stripe non configuré")
