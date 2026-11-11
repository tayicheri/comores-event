package mail

import (
	"strings"
	"testing"

	"comores-event-api/internal/models"
)

func TestBuildTicketConfirmation(t *testing.T) {
	b := &models.Billet{
		ID:          "abc-123",
		FirstName:   "Amina",
		LastName:    "Said",
		Email:       "a@example.com",
		City:        "Moroni",
		Quantity:    2,
		DayVendredi: true,
		DaySamedi:   true,
		AmountCents: 6000,
	}
	subject, textBody, htmlBody := buildTicketConfirmation(b)
	if !strings.Contains(subject, "Confirmation") {
		t.Fatalf("subject: %q", subject)
	}
	for _, needle := range []string{"abc-123", "Amina Said", "2", "60.00 €", "Vendredi et Samedi"} {
		if !strings.Contains(textBody, needle) {
			t.Fatalf("text missing %q in %q", needle, textBody)
		}
		if !strings.Contains(htmlBody, needle) {
			t.Fatalf("html missing %q", needle)
		}
	}
}

func TestFormatDays(t *testing.T) {
	if got := formatDays(true, false); got != "Vendredi" {
		t.Fatalf("got %q", got)
	}
	if got := formatDays(false, true); got != "Samedi" {
		t.Fatalf("got %q", got)
	}
	if got := formatDays(true, true); !strings.Contains(got, "pass") {
		t.Fatalf("got %q", got)
	}
}

func TestSanitizeHeader(t *testing.T) {
	got := sanitizeHeader("ok\r\ninject")
	if strings.Contains(got, "\r") || strings.Contains(got, "\n") {
		t.Fatalf("header injection not stripped: %q", got)
	}
}
