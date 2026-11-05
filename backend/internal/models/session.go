package models

import "time"

// RefreshToken correspond à la table refresh_tokens.
// Seul le hash SHA-256 du token opaque est stocké.
type RefreshToken struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	SessionID  string     `db:"session_id"`
	TokenHash  string     `db:"token_hash"`
	DeviceName *string    `db:"device_name"`
	DeviceOS   *string    `db:"device_os"`
	IPAddress  *string    `db:"ip_address"`
	UserAgent  *string    `db:"user_agent"`
	LastSeenAt time.Time  `db:"last_seen_at"`
	ExpiresAt  time.Time  `db:"expires_at"`
	RevokedAt  *time.Time `db:"revoked_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

// SessionMeta décrit l'appareil à l'origine d'une session.
type SessionMeta struct {
	DeviceName string
	DeviceOS   string
	IPAddress  string
	UserAgent  string
}
