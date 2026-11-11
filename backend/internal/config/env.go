// Package config charge la configuration depuis l'environnement.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config contient toute la configuration de l'API.
type Config struct {
	AppEnv string
	Port   string
	TZ     string

	JWTSecret     string
	FrontendURL   string
	CookieDomain  string
	TrustedProxy  string
	RedisURL      string
	RateLimitCfg  RateLimits
	DB            DBConfig
	CookieSecure  bool
	IsProduction  bool
	AllowedOrigin []string

	StripeSecretKey      string
	StripeWebhookSecret  string
	TicketPriceDayCents  int
	TicketPricePassCents int

	SMTP SMTPConfig
}

// SMTPConfig : envoi des emails client (confirmation billet), distinct du SMTP Grafana.
type SMTPConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	FromAddress string
	FromName    string
}

// Configured indique si l'envoi SMTP est utilisable.
func (s SMTPConfig) Configured() bool {
	return strings.TrimSpace(s.Host) != "" && strings.TrimSpace(s.FromAddress) != ""
}

// RateLimits : nombre de requêtes autorisées par minute et par IP.
type RateLimits struct {
	Login     int
	Register  int
	Refresh   int
	APIGlobal int
	Payments  int
}

// DBConfig : paramètres PostgreSQL.
type DBConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

// Load lit .env (s'il existe) puis l'environnement.
func Load() (*Config, error) {
	_ = godotenv.Load()

	env := strings.ToLower(getEnv("APP_ENV", "development"))
	cfg := &Config{
		AppEnv:       env,
		IsProduction: env == "production",
		Port:         getEnv("PORT", "8080"),
		TZ:           getEnv("TZ", "UTC"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
		CookieDomain: strings.TrimSpace(os.Getenv("COOKIE_DOMAIN")),
		TrustedProxy: getEnv("TRUSTED_PROXIES", "false"),
		RedisURL:     strings.TrimSpace(os.Getenv("REDIS_URL")),
		RateLimitCfg: RateLimits{
			Login:     getEnvInt("RATE_LIMIT_LOGIN", 10),
			Register:  getEnvInt("RATE_LIMIT_REGISTER", 5),
			Refresh:   getEnvInt("RATE_LIMIT_REFRESH", 30),
			APIGlobal: getEnvInt("RATE_LIMIT_API_GLOBAL", 300),
			Payments:  getEnvInt("RATE_LIMIT_PAYMENTS", 20),
		},
		DB: DBConfig{
			Host:    firstEnv([]string{"DB_HOST_POSTGRES", "DB_HOST"}, "db-postgres"),
			Port:    getEnv("DB_PORT", "5432"),
			User:    getEnv("DB_USER", "postgres"),
			Pass:    firstEnv([]string{"DB_PASS", "DB_PASSWORD"}, ""),
			Name:    getEnv("DB_NAME", "postgres"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		StripeSecretKey:      strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY")),
		StripeWebhookSecret:  strings.TrimSpace(os.Getenv("STRIPE_WEBHOOK_SECRET")),
		TicketPriceDayCents:  getEnvInt("TICKET_PRICE_DAY_CENTS", 2000),
		TicketPricePassCents: getEnvInt("TICKET_PRICE_PASS_CENTS", 3000),
		SMTP: SMTPConfig{
			Host:        strings.TrimSpace(os.Getenv("SMTP_HOST")),
			Port:        getEnv("SMTP_PORT", "587"),
			User:        strings.TrimSpace(os.Getenv("SMTP_USER")),
			Password:    os.Getenv("SMTP_PASSWORD"),
			FromAddress: strings.TrimSpace(os.Getenv("SMTP_FROM_ADDRESS")),
			FromName:    getEnv("SMTP_FROM_NAME", "Comores Event"),
		},
	}
	cfg.CookieSecure = cfg.IsProduction
	cfg.AllowedOrigin = splitOrigins(cfg.FrontendURL)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET est obligatoire")
	}
	if c.IsProduction && len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET doit contenir au moins 32 caractères en production (actuel : %d)", len(c.JWTSecret))
	}
	if c.TicketPriceDayCents <= 0 || c.TicketPricePassCents <= 0 {
		return errors.New("TICKET_PRICE_DAY_CENTS et TICKET_PRICE_PASS_CENTS doivent être > 0")
	}
	if c.IsProduction && c.StripeSecretKey == "" {
		return errors.New("STRIPE_SECRET_KEY est obligatoire en production")
	}
	return nil
}

func splitOrigins(raw string) []string {
	var out []string
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimRight(strings.TrimSpace(o), "/")
		if o != "" {
			out = append(out, o)
		}
	}
	return out
}

func getEnv(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func firstEnv(keys []string, def string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
