package middleware

import (
	"net/http"
	"time"

	"comores-event-api/internal/constants"
)

// CookieConfig regroupe les attributs communs des cookies de session.
type CookieConfig struct {
	Secure bool   // true quand APP_ENV=production
	Domain string // COOKIE_DOMAIN (vide = host-only)
}

func (c CookieConfig) base(name, value, path string, ttl time.Duration, httpOnly bool) *http.Cookie {
	ck := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   c.Domain,
		HttpOnly: httpOnly,
		Secure:   c.Secure,
		SameSite: http.SameSiteLaxMode,
	}
	if ttl > 0 {
		ck.MaxAge = int(ttl.Seconds())
		ck.Expires = time.Now().Add(ttl)
	} else {
		ck.MaxAge = -1
		ck.Expires = time.Unix(0, 0)
	}
	return ck
}

// SetSessionCookies pose les trois cookies : accès (HttpOnly), refresh (HttpOnly), CSRF (lisible JS).
func (c CookieConfig) SetSessionCookies(w http.ResponseWriter, accessJWT, refreshToken, csrfToken string) {
	http.SetCookie(w, c.base(constants.CookieAccessToken, accessJWT, "/", constants.AccessTokenTTL, true))
	http.SetCookie(w, c.base(constants.CookieRefreshToken, refreshToken, constants.RefreshCookiePath, constants.RefreshTokenTTL, true))
	http.SetCookie(w, c.base(constants.CookieCSRF, csrfToken, "/", constants.CSRFCookieTTL, false))
}

// ClearSessionCookies expire les trois cookies.
func (c CookieConfig) ClearSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, c.base(constants.CookieAccessToken, "", "/", 0, true))
	http.SetCookie(w, c.base(constants.CookieRefreshToken, "", constants.RefreshCookiePath, 0, true))
	http.SetCookie(w, c.base(constants.CookieCSRF, "", "/", 0, false))
}
