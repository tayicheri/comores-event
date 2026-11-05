package payments

import (
	"fmt"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
)

// StripeAPI abstrait les appels Stripe (mockable en tests).
type StripeAPI interface {
	CreatePaymentIntent(amount int, currency, receiptEmail string, metadata map[string]string) (id, clientSecret, status string, err error)
	RetrievePaymentIntent(id string) (status string, amount int, currency string, metadata map[string]string, err error)
}

// LiveStripe implémentation réelle via stripe-go.
type LiveStripe struct {
	SecretKey string
}

func (s LiveStripe) key() (string, error) {
	key := strings.TrimSpace(s.SecretKey)
	if key == "" {
		return "", fmt.Errorf("STRIPE_SECRET_KEY manquante")
	}
	return key, nil
}

func (s LiveStripe) CreatePaymentIntent(amount int, currency, receiptEmail string, metadata map[string]string) (string, string, string, error) {
	key, err := s.key()
	if err != nil {
		return "", "", "", err
	}
	stripe.Key = key
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amount)),
		Currency:           stripe.String(strings.ToLower(currency)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}
	if receiptEmail != "" {
		params.ReceiptEmail = stripe.String(receiptEmail)
	}
	for k, v := range metadata {
		params.AddMetadata(k, v)
	}
	intent, err := paymentintent.New(params)
	if err != nil {
		return "", "", "", err
	}
	return intent.ID, intent.ClientSecret, string(intent.Status), nil
}

func (s LiveStripe) RetrievePaymentIntent(id string) (string, int, string, map[string]string, error) {
	key, err := s.key()
	if err != nil {
		return "", 0, "", nil, err
	}
	stripe.Key = key
	intent, err := paymentintent.Get(id, nil)
	if err != nil {
		return "", 0, "", nil, err
	}
	meta := map[string]string{}
	for k, v := range intent.Metadata {
		meta[k] = v
	}
	return string(intent.Status), int(intent.Amount), string(intent.Currency), meta, nil
}

// ConstructStripeEvent vérifie la signature webhook Stripe.
func ConstructStripeEvent(payload []byte, header, secret string) (stripe.Event, error) {
	return webhook.ConstructEventWithOptions(payload, header, secret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
}
