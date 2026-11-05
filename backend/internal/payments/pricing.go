// Package payments : encaissement Stripe et calcul des montants billets.
package payments

import (
	"fmt"

	"comores-event-api/internal/constants"
)

// PricingConfig prix unitaires serveur (centimes EUR).
type PricingConfig struct {
	DayCents  int
	PassCents int
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

// AmountCents recalcule le montant total côté serveur.
func AmountCents(cfg PricingConfig, formula string, quantity int) (int, error) {
	if quantity < 1 || quantity > constants.MaxTicketQuantity {
		return 0, fmt.Errorf("quantité invalide")
	}
	var unit int
	switch formula {
	case constants.TicketFormulaDay:
		unit = cfg.DayCents
	case constants.TicketFormulaPass:
		unit = cfg.PassCents
	default:
		return 0, fmt.Errorf("formule invalide")
	}
	if unit <= 0 {
		return 0, fmt.Errorf("tarif non configuré")
	}
	return unit * quantity, nil
}
