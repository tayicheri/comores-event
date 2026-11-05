// Package constants regroupe les constantes d'authentification partagées.
package constants

import "time"

const (
	// Cookies
	CookieAccessToken  = "auth_token"
	CookieRefreshToken = "refresh_token"
	CookieCSRF         = "XSRF-TOKEN"
	HeaderCSRF         = "X-XSRF-TOKEN"

	// Chemin du cookie de refresh (limité aux routes d'auth).
	RefreshCookiePath = "/api/auth"

	// Durées de vie
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
	CSRFCookieTTL   = RefreshTokenTTL

	// Sessions
	MaxActiveSessions = 5

	// Fenêtre de tolérance pour deux refresh concurrents avec le même token.
	RefreshReuseGrace = 10 * time.Second

	// Verrouillage de compte
	MaxFailedLogins = 5
	LockDuration    = 15 * time.Minute

	// JWT
	JWTIssuer = "comores-event-api"

	// Mots de passe (bcrypt tronque / refuse au-delà de 72 octets)
	PasswordMinLen = 8
	PasswordMaxLen = 72

	// Rôles
	RoleUser  = "user"
	RoleAdmin = "admin"

	// Statuts de compte
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusPending  = "pending"
)

// CSRFExemptPaths liste les routes (POST) exemptées de la vérification CSRF.
var CSRFExemptPaths = map[string]struct{}{
	"/api/auth/login":    {},
	"/api/auth/register": {},
	"/api/auth/refresh":  {},
}
