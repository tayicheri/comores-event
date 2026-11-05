package api

import (
	"net/http"

	"comores-event-api/internal/handlers"
	"comores-event-api/internal/middleware"
)

// registerAuthRoutes enregistre /api/auth/*.
// Chaque route /api/ est soumise à la limite globale (RATE_LIMIT_API_GLOBAL) ;
// login/register/refresh ont en plus leur limite dédiée.
func registerAuthRoutes(mux *http.ServeMux, d Deps) {
	cookies := middleware.CookieConfig{Secure: d.Config.CookieSecure, Domain: d.Config.CookieDomain}
	h := handlers.NewAuthHandler(d.Users, d.Sessions, d.JWT, cookies, d.Proxies)

	rl := d.Config.RateLimitCfg
	global := middleware.RateLimit(d.Limiter, d.Proxies, "api", rl.APIGlobal, rateWindow)
	loginRL := middleware.RateLimit(d.Limiter, d.Proxies, "login", rl.Login, rateWindow)
	registerRL := middleware.RateLimit(d.Limiter, d.Proxies, "register", rl.Register, rateWindow)
	refreshRL := middleware.RateLimit(d.Limiter, d.Proxies, "refresh", rl.Refresh, rateWindow)
	requireAuth := middleware.Auth(d.JWT, d.Users)

	route := func(pattern string, fn http.HandlerFunc, mws ...middleware.Middleware) {
		mux.Handle(pattern, middleware.Chain(fn, append([]middleware.Middleware{global}, mws...)...))
	}

	route("POST /api/auth/login", h.Login, loginRL)
	route("POST /api/auth/register", h.Register, registerRL)
	route("POST /api/auth/refresh", h.Refresh, refreshRL)
	route("POST /api/auth/logout", h.Logout)
	route("POST /api/auth/logout-all", h.LogoutAll, requireAuth)
	route("GET /api/auth/me", h.Me, requireAuth)
	route("PATCH /api/auth/me", h.UpdateMe, requireAuth)
}
