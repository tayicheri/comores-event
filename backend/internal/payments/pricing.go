// Package payments : encaissement Stripe et calcul des montants billets.
package payments

import (
	"fmt"
	"time"

	"comores-event-api/internal/constants"
)

// PassPresaleEndsAt : fin de la prévente pass (Europe/Paris). À partir de cet instant → tarif régulier.
var PassPresaleEndsAt = time.Date(2027, 1, 1, 0, 0, 0, 0, mustLoadLocation("Europe/Paris"))

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic("payments: fuseau horaire introuvable: " + name)
	}
	return loc
}

// PricingConfig prix unitaires serveur (centimes EUR).
type PricingConfig struct {
	DayCents         int
	PassPresaleCents int
	PassRegularCents int
}

// IsPassPresale indique si now est strictement avant la fin de prévente.
func IsPassPresale(now time.Time) bool {
	return now.Before(PassPresaleEndsAt)
}

// PassUnitCents retourne le tarif unitaire pass selon la date.
func PassUnitCents(cfg PricingConfig, now time.Time) int {
	if IsPassPresale(now) {
		return cfg.PassPresaleCents
	}
	return cfg.PassRegularCents
}

// DayFlags calcule les booléens jour à partir de la formule / jour choisi.
func DayFlags(formula, day string) (vendredi, samedi bool, err error) {
	switch formula {
	case constants.TicketFormulaPass:
		return true, true, nil
	case constants.TicketFormulaDay:
		switch day {
		case constants.TicketDayVendredi:
			return true, false, nil
		case constants.TicketDaySamedi:
			return false, true, nil
		default:
			return false, false, fmt.Errorf("jour requis pour un billet journée (vendredi ou samedi)")
		}
	default:
		return false, false, fmt.Errorf("formule invalide")
	}
}

// AmountCents recalcule le montant total côté serveur (horloge serveur).
func AmountCents(cfg PricingConfig, formula string, quantity int) (int, error) {
	return AmountCentsAt(cfg, formula, quantity, time.Now())
}

// AmountCentsAt recalcule le montant total à une date donnée (tests / bascule prévente).
func AmountCentsAt(cfg PricingConfig, formula string, quantity int, now time.Time) (int, error) {
	if quantity < 1 || quantity > constants.MaxTicketQuantity {
		return 0, fmt.Errorf("quantité invalide")
	}
	var unit int
	switch formula {
	case constants.TicketFormulaDay:
		unit = cfg.DayCents
	case constants.TicketFormulaPass:
		unit = PassUnitCents(cfg, now)
	default:
		return 0, fmt.Errorf("formule invalide")
	}
	if unit <= 0 {
		return 0, fmt.Errorf("tarif non configuré")
	}
	return unit * quantity, nil
}
