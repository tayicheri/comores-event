package payments

import (
	"testing"
	"time"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/models"
)

func TestDayFlags(t *testing.T) {
	v, s, err := DayFlags(constants.TicketFormulaPass, "")
	if err != nil || !v || !s {
		t.Fatalf("pass: v=%v s=%v err=%v", v, s, err)
	}
	v, s, err = DayFlags(constants.TicketFormulaDay, constants.TicketDayVendredi)
	if err != nil || !v || s {
		t.Fatalf("vendredi: v=%v s=%v err=%v", v, s, err)
	}
	v, s, err = DayFlags(constants.TicketFormulaDay, constants.TicketDaySamedi)
	if err != nil || v || !s {
		t.Fatalf("samedi: v=%v s=%v err=%v", v, s, err)
	}
	if _, _, err := DayFlags(constants.TicketFormulaDay, ""); err == nil {
		t.Fatal("day without day should fail")
	}
}

func TestAmountCents(t *testing.T) {
	cfg := PricingConfig{DayCents: 2000, PassPresaleCents: 3000, PassRegularCents: 3500}
	got, err := AmountCentsAt(cfg, constants.TicketFormulaDay, 3, time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))
	if err != nil || got != 6000 {
		t.Fatalf("day: got=%d err=%v", got, err)
	}
	got, err = AmountCentsAt(cfg, constants.TicketFormulaPass, 2, time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))
	if err != nil || got != 6000 {
		t.Fatalf("pass presale: got=%d err=%v", got, err)
	}
	if _, err := AmountCentsAt(cfg, constants.TicketFormulaDay, 0, time.Now()); err == nil {
		t.Fatal("qty 0 should fail")
	}
}

func TestPassUnitCentsPresaleCutoff(t *testing.T) {
	cfg := PricingConfig{DayCents: 2000, PassPresaleCents: 3000, PassRegularCents: 3500}
	paris := PassPresaleEndsAt.Location()

	before := PassPresaleEndsAt.Add(-time.Second)
	if got := PassUnitCents(cfg, before); got != 3000 {
		t.Fatalf("before cutoff: got=%d want=3000", got)
	}
	if !IsPassPresale(before) {
		t.Fatal("before cutoff should be presale")
	}

	at := time.Date(2027, 1, 1, 0, 0, 0, 0, paris)
	if got := PassUnitCents(cfg, at); got != 3500 {
		t.Fatalf("at cutoff: got=%d want=3500", got)
	}
	if IsPassPresale(at) {
		t.Fatal("at cutoff should not be presale")
	}

	after := PassPresaleEndsAt.Add(time.Hour)
	if got := PassUnitCents(cfg, after); got != 3500 {
		t.Fatalf("after cutoff: got=%d want=3500", got)
	}

	got, err := AmountCentsAt(cfg, constants.TicketFormulaPass, 1, at)
	if err != nil || got != 3500 {
		t.Fatalf("amount at cutoff: got=%d err=%v", got, err)
	}
}

func TestInputFromMetadataRoundTrip(t *testing.T) {
	in := models.CreateBilletInput{
		FirstName:             "Amina",
		LastName:              "Said",
		Email:                 "amina@example.com",
		City:                  "Moroni",
		Quantity:              2,
		DayVendredi:           true,
		DaySamedi:             false,
		AmountCents:           4000,
		StripePaymentIntentID: "pi_test",
	}
	meta := BuildMetadata(in)
	out, err := InputFromMetadata("pi_test", meta)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("roundtrip mismatch\n%+v\n%+v", out, in)
	}
}
