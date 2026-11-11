// Package mail envoie les emails transactionnels de l'API (confirmation billet).
package mail

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strings"
	"sync"

	"comores-event-api/internal/config"
	"comores-event-api/internal/models"
)

// Mailer envoie des emails métier.
type Mailer interface {
	SendTicketConfirmation(ctx context.Context, b *models.Billet) error
}

// NopMailer ne fait rien (SMTP non configuré / tests).
type NopMailer struct{}

// SendTicketConfirmation est un no-op.
func (NopMailer) SendTicketConfirmation(context.Context, *models.Billet) error {
	return nil
}

// SMTPMailer envoie via SMTP (STARTTLS si annoncé par le serveur).
type SMTPMailer struct {
	cfg config.SMTPConfig
}

var warnSMTPOnce sync.Once

// NewFromConfig retourne un SMTPMailer si configuré, sinon NopMailer.
func NewFromConfig(cfg config.SMTPConfig) Mailer {
	if !cfg.Configured() {
		warnSMTPOnce.Do(func() {
			slog.Warn("SMTP non configuré : emails de confirmation billet désactivés")
		})
		return NopMailer{}
	}
	if strings.TrimSpace(cfg.Port) == "" {
		cfg.Port = "587"
	}
	if strings.TrimSpace(cfg.FromName) == "" {
		cfg.FromName = "Comores Event"
	}
	return &SMTPMailer{cfg: cfg}
}

// SendTicketConfirmation envoie la confirmation d'émission au titulaire du billet.
func (m *SMTPMailer) SendTicketConfirmation(ctx context.Context, b *models.Billet) error {
	if b == nil {
		return fmt.Errorf("billet nil")
	}
	to := strings.TrimSpace(b.Email)
	if to == "" {
		return fmt.Errorf("email destinataire vide")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	subject, textBody, htmlBody := buildTicketConfirmation(b)
	msg := buildMIMEMessage(m.cfg.FromName, m.cfg.FromAddress, to, subject, textBody, htmlBody)

	addr := net.JoinHostPort(m.cfg.Host, m.cfg.Port)
	var auth smtp.Auth
	if m.cfg.User != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
	}
	if err := smtp.SendMail(addr, auth, m.cfg.FromAddress, []string{to}, msg); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}
