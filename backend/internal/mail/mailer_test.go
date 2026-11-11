package mail

import (
	"testing"

	"comores-event-api/internal/config"
)

func TestNewFromConfigNopWhenIncomplete(t *testing.T) {
	m := NewFromConfig(config.SMTPConfig{})
	if _, ok := m.(NopMailer); !ok {
		t.Fatalf("want NopMailer, got %T", m)
	}
	m = NewFromConfig(config.SMTPConfig{Host: "smtp.example.com"})
	if _, ok := m.(NopMailer); !ok {
		t.Fatalf("want NopMailer without FromAddress, got %T", m)
	}
}

func TestNewFromConfigSMTPWhenReady(t *testing.T) {
	m := NewFromConfig(config.SMTPConfig{
		Host:        "smtp.example.com",
		Port:        "587",
		FromAddress: "billets@example.com",
		FromName:    "Comores Event",
	})
	if _, ok := m.(*SMTPMailer); !ok {
		t.Fatalf("want *SMTPMailer, got %T", m)
	}
}
