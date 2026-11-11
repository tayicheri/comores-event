package mail

import (
	"fmt"
	"html"
	"strings"

	"comores-event-api/internal/models"
)

func buildTicketConfirmation(b *models.Billet) (subject, textBody, htmlBody string) {
	days := formatDays(b.DayVendredi, b.DaySamedi)
	amount := formatEUR(b.AmountCents)
	name := strings.TrimSpace(b.FirstName + " " + b.LastName)

	subject = "Confirmation de votre billet — Comores Event"

	textBody = fmt.Sprintf(`Bonjour %s,

Votre billet a bien été émis. Voici le récapitulatif :

- Référence : %s
- Jours : %s
- Quantité : %d
- Montant : %s
- Ville : %s

Présentez-vous à l'entrée avec une pièce d'identité et cette confirmation.

À bientôt,
L'équipe Comores Event
`, name, b.ID, days, b.Quantity, amount, b.City)

	htmlBody = fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head><meta charset="utf-8"><title>%s</title></head>
<body style="font-family:system-ui,sans-serif;line-height:1.5;color:#111">
  <p>Bonjour %s,</p>
  <p>Votre billet a bien été émis. Voici le récapitulatif :</p>
  <ul>
    <li><strong>Référence :</strong> %s</li>
    <li><strong>Jours :</strong> %s</li>
    <li><strong>Quantité :</strong> %d</li>
    <li><strong>Montant :</strong> %s</li>
    <li><strong>Ville :</strong> %s</li>
  </ul>
  <p>Présentez-vous à l'entrée avec une pièce d'identité et cette confirmation.</p>
  <p>À bientôt,<br>L'équipe Comores Event</p>
</body>
</html>
`,
		html.EscapeString(subject),
		html.EscapeString(name),
		html.EscapeString(b.ID),
		html.EscapeString(days),
		b.Quantity,
		html.EscapeString(amount),
		html.EscapeString(b.City),
	)
	return subject, textBody, htmlBody
}

func formatDays(vendredi, samedi bool) string {
	switch {
	case vendredi && samedi:
		return "Vendredi et Samedi (pass)"
	case vendredi:
		return "Vendredi"
	case samedi:
		return "Samedi"
	default:
		return "—"
	}
}

func formatEUR(cents int) string {
	euros := float64(cents) / 100
	return fmt.Sprintf("%.2f €", euros)
}

func buildMIMEMessage(fromName, fromAddr, to, subject, textBody, htmlBody string) []byte {
	boundary := "comores-event-boundary"
	fromHeader := fromAddr
	if strings.TrimSpace(fromName) != "" {
		fromHeader = fmt.Sprintf("%s <%s>", sanitizeHeader(fromName), fromAddr)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
	b.WriteString(fmt.Sprintf("To: %s\r\n", to))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", sanitizeHeader(subject)))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundary))
	b.WriteString("\r\n")
	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(textBody)
	b.WriteString("\r\n")
	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(htmlBody)
	b.WriteString("\r\n")
	b.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	return []byte(b.String())
}

func sanitizeHeader(v string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		}
		return r
	}, v)
}
