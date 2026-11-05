package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"comores-event-api/internal/auth"
	"comores-event-api/internal/constants"
	"comores-event-api/internal/middleware"
	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

// sessionIssuer émet les sessions (refresh token + JWT + cookies).
type sessionIssuer struct {
	sessions repository.SessionRepository
	jwt      *auth.JWTManager
	cookies  middleware.CookieConfig
	proxies  *middleware.TrustedProxies
}

// issueNew crée une nouvelle session (login/register) : applique la limite de
// sessions actives, insère le refresh token et pose les cookies.
func (s *sessionIssuer) issueNew(ctx context.Context, w http.ResponseWriter, r *http.Request, user *models.User) error {
	if err := s.sessions.EnforceSessionLimit(ctx, user.ID, constants.MaxActiveSessions-1); err != nil {
		return err
	}
	raw, tok, err := s.newRefreshToken(r)
	if err != nil {
		return err
	}
	tok.UserID = user.ID
	if err := s.sessions.Create(ctx, tok); err != nil {
		return err
	}
	return s.writeCookies(w, user, tok.SessionID, raw)
}

// rotate remplace le refresh token oldRaw par un nouveau (même session) et pose les cookies.
// loadUser est appelé après rotation pour récupérer l'utilisateur propriétaire.
func (s *sessionIssuer) rotate(ctx context.Context, w http.ResponseWriter, r *http.Request, oldRaw string,
	loadUser func(ctx context.Context, userID string) (*models.User, error)) (*models.User, error) {

	raw, tok, err := s.newRefreshToken(r)
	if err != nil {
		return nil, err
	}
	if err := s.sessions.Rotate(ctx, auth.HashToken(oldRaw), tok); err != nil {
		return nil, err
	}
	user, err := loadUser(ctx, tok.UserID)
	if err == nil && !user.IsActive() {
		err = repository.ErrNotFound
	}
	if err != nil {
		_ = s.sessions.RevokeToken(ctx, tok.ID)
		return nil, err
	}
	if err := s.writeCookies(w, user, tok.SessionID, raw); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *sessionIssuer) newRefreshToken(r *http.Request) (string, *models.RefreshToken, error) {
	raw, err := auth.GenerateOpaqueToken()
	if err != nil {
		return "", nil, err
	}
	meta := s.sessionMeta(r)
	tok := &models.RefreshToken{
		TokenHash: auth.HashToken(raw),
		ExpiresAt: time.Now().Add(constants.RefreshTokenTTL),
	}
	tok.DeviceName = strPtr(meta.DeviceName)
	tok.DeviceOS = strPtr(meta.DeviceOS)
	tok.IPAddress = strPtr(meta.IPAddress)
	tok.UserAgent = strPtr(meta.UserAgent)
	return raw, tok, nil
}

func (s *sessionIssuer) writeCookies(w http.ResponseWriter, user *models.User, sessionID, rawRefresh string) error {
	access, err := s.jwt.Generate(user.ID, user.PrimaryRole(), user.TokenVersion, sessionID)
	if err != nil {
		return err
	}
	csrf, err := auth.GenerateCSRFToken()
	if err != nil {
		return err
	}
	s.cookies.SetSessionCookies(w, access, rawRefresh, csrf)
	return nil
}

func (s *sessionIssuer) sessionMeta(r *http.Request) models.SessionMeta {
	ua := r.UserAgent()
	if len(ua) > 512 {
		ua = ua[:512]
	}
	name, os := parseDevice(ua)
	return models.SessionMeta{
		DeviceName: name,
		DeviceOS:   os,
		IPAddress:  s.proxies.ClientIP(r),
		UserAgent:  ua,
	}
}

// parseDevice : détection grossière navigateur / OS depuis le User-Agent.
func parseDevice(ua string) (name, os string) {
	l := strings.ToLower(ua)
	switch {
	case ua == "":
		return "", ""
	case strings.Contains(l, "edg/"):
		name = "Edge"
	case strings.Contains(l, "opr/") || strings.Contains(l, "opera"):
		name = "Opera"
	case strings.Contains(l, "firefox"):
		name = "Firefox"
	case strings.Contains(l, "chrome") || strings.Contains(l, "crios"):
		name = "Chrome"
	case strings.Contains(l, "safari"):
		name = "Safari"
	default:
		name = "Autre"
	}
	switch {
	case strings.Contains(l, "android"):
		os = "Android"
	case strings.Contains(l, "iphone") || strings.Contains(l, "ipad") || strings.Contains(l, "ios"):
		os = "iOS"
	case strings.Contains(l, "windows"):
		os = "Windows"
	case strings.Contains(l, "mac os") || strings.Contains(l, "macintosh"):
		os = "macOS"
	case strings.Contains(l, "linux"):
		os = "Linux"
	default:
		os = "Autre"
	}
	return name, os
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
