package payments

import (
	"testing"

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
	cfg := PricingConfig{DayCents: 2000, PassCents: 3000}
	got, err := AmountCents(cfg, constants.TicketFormulaDay, 3)
	if err != nil || got != 6000 {
		t.Fatalf("day: got=%d err=%v", got, err)
	}
	got, err = AmountCents(cfg, constants.TicketFormulaPass, 2)
	if err != nil || got != 6000 {
		t.Fatalf("pass: got=%d err=%v", got, err)
	}
	if _, err := AmountCents(cfg, constants.TicketFormulaDay, 0); err == nil {
		t.Fatal("qty 0 should fail")
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
		t.Fatalf("roundtrip mismatch\n%+v\n%+v", in, out)
	}
}
