package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"comores-event-api/internal/auth"
	"comores-event-api/internal/constants"
	"comores-event-api/internal/middleware"
	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

// AuthHandler regroupe les handlers d'authentification.
type AuthHandler struct {
	users   repository.UserRepository
	session *sessionIssuer
	cookies middleware.CookieConfig
}

// NewAuthHandler construit le handler d'authentification.
func NewAuthHandler(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	jwtm *auth.JWTManager,
	cookies middleware.CookieConfig,
	proxies *middleware.TrustedProxies,
) *AuthHandler {
	return &AuthHandler{
		users:   users,
		cookies: cookies,
		session: &sessionIssuer{sessions: sessions, jwt: jwtm, cookies: cookies, proxies: proxies},
	}
}

func normalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

func internalError(w http.ResponseWriter, msg string, err error) {
	slog.Error(msg, "error", err)
	WriteError(w, http.StatusInternalServerError, "Erreur interne du serveur")
}

// Register : POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Email = normalizeEmail(req.Email)
	if req.DisplayName != nil {
		req.DisplayName = strPtr(strings.TrimSpace(*req.DisplayName))
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		internalError(w, "register: hash", err)
		return
	}
	user, err := h.users.Create(r.Context(), req.Email, hash, req.DisplayName)
	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			WriteError(w, http.StatusConflict, "Cette adresse email est déjà utilisée")
			return
		}
		internalError(w, "register: create", err)
		return
	}
	if err := h.session.issueNew(r.Context(), w, r, user); err != nil {
		internalError(w, "register: session", err)
		return
	}
	WriteJSON(w, http.StatusCreated, user.ToResponse())
}

// Login : POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	email := normalizeEmail(req.Email)
	const invalid = "Email ou mot de passe incorrect"

	user, err := h.users.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			auth.BurnPasswordCheck(req.Password)
			WriteError(w, http.StatusUnauthorized, invalid)
			return
		}
		internalError(w, "login: lookup", err)
		return
	}
	if user.DeletedAt != nil {
		auth.BurnPasswordCheck(req.Password)
		WriteError(w, http.StatusUnauthorized, invalid)
		return
	}
	if user.IsLocked(time.Now()) {
		WriteError(w, http.StatusTooManyRequests, "Compte temporairement verrouillé, réessayez plus tard")
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		if err := h.users.RecordFailedLogin(r.Context(), user.ID); err != nil {
			slog.Error("login: record failed", "error", err)
		}
		WriteError(w, http.StatusUnauthorized, invalid)
		return
	}
	switch user.AccountStatus {
	case constants.StatusActive:
	case constants.StatusPending:
		WriteError(w, http.StatusForbidden, "Compte en attente d'activation")
		return
	default:
		WriteError(w, http.StatusForbidden, "Compte désactivé")
		return
	}

	if err := h.users.ResetFailedLogins(r.Context(), user.ID); err != nil {
		slog.Error("login: reset failed", "error", err)
	}
	if err := h.session.issueNew(r.Context(), w, r, user); err != nil {
		internalError(w, "login: session", err)
		return
	}
	WriteJSON(w, http.StatusOK, user.ToResponse())
}

// Refresh : POST /api/auth/refresh — rotation du refresh token.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(constants.CookieRefreshToken)
	if err != nil || cookie.Value == "" {
		h.cookies.ClearSessionCookies(w)
		WriteError(w, http.StatusUnauthorized, "Session expirée")
		return
	}

	user, err := h.session.rotate(r.Context(), w, r, cookie.Value, h.users.GetByID)
	if err != nil {
		if errors.Is(err, repository.ErrTokenInvalid) ||
			errors.Is(err, repository.ErrTokenReused) ||
			errors.Is(err, repository.ErrNotFound) {
			h.cookies.ClearSessionCookies(w)
			WriteError(w, http.StatusUnauthorized, "Session expirée")
			return
		}
		internalError(w, "refresh", err)
		return
	}
	WriteJSON(w, http.StatusOK, user.ToResponse())
}

// Logout : POST /api/auth/logout — révoque le refresh token courant.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(constants.CookieRefreshToken); err == nil && cookie.Value != "" {
		if err := h.session.sessions.RevokeByHash(r.Context(), auth.HashToken(cookie.Value)); err != nil {
			slog.Error("logout: revoke", "error", err)
		}
	}
	h.cookies.ClearSessionCookies(w)
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Déconnexion réussie"})
}

// LogoutAll : POST /api/auth/logout-all — incrémente token_version et révoque tout.
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentification requise")
		return
	}
	if err := h.users.BumpTokenVersion(r.Context(), user.ID); err != nil {
		internalError(w, "logout-all: bump", err)
		return
	}
	if err := h.session.sessions.RevokeAllForUser(r.Context(), user.ID); err != nil {
		internalError(w, "logout-all: revoke", err)
		return
	}
	h.cookies.ClearSessionCookies(w)
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Toutes les sessions ont été révoquées"})
}

// Me : GET /api/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentification requise")
		return
	}
	WriteJSON(w, http.StatusOK, user.ToResponse())
}

// UpdateMe : PATCH /api/auth/me
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Authentification requise")
		return
	}
	var req models.UpdateProfileRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.DisplayName != nil {
		name := strings.TrimSpace(*req.DisplayName)
		if err := h.users.UpdateDisplayName(r.Context(), user.ID, strPtr(name)); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				WriteError(w, http.StatusNotFound, "Utilisateur introuvable")
				return
			}
			internalError(w, "update me", err)
			return
		}
	}
	updated, err := h.users.GetByID(r.Context(), user.ID)
	if err != nil {
		internalError(w, "update me: reload", err)
		return
	}
	WriteJSON(w, http.StatusOK, updated.ToResponse())
}
