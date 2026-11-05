// Package auth : hachage de mots de passe, JWT d'accès et tokens opaques.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"comores-event-api/internal/constants"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// BcryptCost : coût bcrypt (compatible avec crypt(..., gen_salt('bf', 14)) de pgcrypto).
const BcryptCost = 14

// HashPassword hache un mot de passe avec bcrypt.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword compare un mot de passe avec son hash bcrypt.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

var (
	dummyOnce sync.Once
	dummyHash string
)

// BurnPasswordCheck exécute une comparaison bcrypt factice afin de lisser le
// temps de réponse quand l'utilisateur n'existe pas (anti-énumération).
func BurnPasswordCheck(password string) {
	dummyOnce.Do(func() {
		dummyHash, _ = HashPassword("dummy-password-for-timing")
	})
	_ = CheckPassword(dummyHash, password)
}

// ---------------------------------------------------------------------------
// JWT (HS256)
// ---------------------------------------------------------------------------

// Claims : claims du JWT d'accès.
type Claims struct {
	Role         string `json:"role"`
	TokenVersion int    `json:"tv"`
	SessionID    string `json:"sid"`
	jwt.RegisteredClaims
}

// UserID retourne le sujet (id utilisateur).
func (c *Claims) UserID() string { return c.Subject }

// JWTManager signe et vérifie les JWT d'accès.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTManager crée un gestionnaire JWT HS256.
func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: constants.AccessTokenTTL}
}

// Generate émet un JWT d'accès pour l'utilisateur donné.
func (m *JWTManager) Generate(userID, role string, tokenVersion int, sessionID string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role:         role,
		TokenVersion: tokenVersion,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    constants.JWTIssuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Parse vérifie la signature (HS256 uniquement), l'expiration et l'émetteur.
func (m *JWTManager) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(constants.JWTIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("jwt invalide: %w", err)
	}
	if !tok.Valid || claims.Subject == "" {
		return nil, errors.New("jwt invalide")
	}
	return claims, nil
}

// ---------------------------------------------------------------------------
// Tokens opaques
// ---------------------------------------------------------------------------

// GenerateOpaqueToken retourne un token aléatoire de 32 octets (base64url).
func GenerateOpaqueToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken retourne le SHA-256 hexadécimal d'un token opaque.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// GenerateCSRFToken retourne un token CSRF aléatoire (hex, 32 octets).
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
