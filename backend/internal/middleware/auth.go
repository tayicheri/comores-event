package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"comores-event-api/internal/auth"
	"comores-event-api/internal/constants"
	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

type ctxKey int

const (
	ctxUser ctxKey = iota
	ctxClaims
)

// UserFromContext retourne l'utilisateur authentifié (chargé avec ses rôles).
func UserFromContext(ctx context.Context) (*models.User, bool) {
	u, ok := ctx.Value(ctxUser).(*models.User)
	return u, ok && u != nil
}

// ClaimsFromContext retourne les claims JWT de la requête.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	c, ok := ctx.Value(ctxClaims).(*auth.Claims)
	return c, ok && c != nil
}

// Auth exige un JWT valide dans le cookie auth_token, vérifie token_version,
// l'état du compte, charge les rôles puis place l'utilisateur dans le contexte.
func Auth(jwtm *auth.JWTManager, users repository.UserRepository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(constants.CookieAccessToken)
			if err != nil || cookie.Value == "" {
				writeError(w, http.StatusUnauthorized, "Authentification requise")
				return
			}
			claims, err := jwtm.Parse(cookie.Value)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Session invalide ou expirée")
				return
			}

			user, err := users.GetByID(r.Context(), claims.UserID())
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					writeError(w, http.StatusUnauthorized, "Session invalide ou expirée")
					return
				}
				slog.Error("auth: chargement utilisateur", "error", err)
				writeError(w, http.StatusInternalServerError, "Erreur interne du serveur")
				return
			}
			if !user.IsActive() || user.TokenVersion != claims.TokenVersion {
				writeError(w, http.StatusUnauthorized, "Session invalide ou expirée")
				return
			}

			ctx := context.WithValue(r.Context(), ctxUser, user)
			ctx = context.WithValue(ctx, ctxClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole exige que l'utilisateur possède l'un des rôles donnés (après Auth).
func RequireRole(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "Authentification requise")
				return
			}
			for _, have := range user.Roles {
				for _, want := range roles {
					if have == want {
						next.ServeHTTP(w, r)
						return
					}
				}
			}
			writeError(w, http.StatusForbidden, "Accès refusé")
		})
	}
}
