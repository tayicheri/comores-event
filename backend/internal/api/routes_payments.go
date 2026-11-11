package api

import (
	"net/http"

	"comores-event-api/internal/constants"
	"comores-event-api/internal/handlers"
	"comores-event-api/internal/mail"
	"comores-event-api/internal/middleware"
	"comores-event-api/internal/payments"
)

func registerPaymentRoutes(mux *http.ServeMux, d Deps) {
	stripeClient := payments.LiveStripe{SecretKey: d.Config.StripeSecretKey}
	pricing := payments.PricingConfig{
		DayCents:  d.Config.TicketPriceDayCents,
		PassCents: d.Config.TicketPricePassCents,
	}
	mailer := mail.NewFromConfig(d.Config.SMTP)
	h := handlers.NewPaymentHandler(d.Billets, stripeClient, pricing, mailer, d.Config.StripeWebhookSecret)

	rl := d.Config.RateLimitCfg
	global := middleware.RateLimit(d.Limiter, d.Proxies, "api", rl.APIGlobal, rateWindow)
	payRL := middleware.RateLimit(d.Limiter, d.Proxies, "payments", rl.Payments, rateWindow)

	route := func(pattern string, fn http.HandlerFunc, mws ...middleware.Middleware) {
		mux.Handle(pattern, middleware.Chain(fn, append([]middleware.Middleware{global}, mws...)...))
	}

	route("POST /api/payments/intent", h.CreateIntent, payRL)
	route("POST /api/payments/confirm", h.Confirm, payRL)
	mux.Handle("POST /api/webhooks/stripe", middleware.Chain(http.HandlerFunc(h.StripeWebhook), global))
}

func registerAdminBilletRoutes(mux *http.ServeMux, d Deps) {
	h := handlers.NewBilletAdminHandler(d.Billets)
	rl := d.Config.RateLimitCfg
	global := middleware.RateLimit(d.Limiter, d.Proxies, "api", rl.APIGlobal, rateWindow)
	auth := middleware.Auth(d.JWT, d.Users)
	role := middleware.RequireRole(constants.RoleAdmin)

	route := func(pattern string, fn http.HandlerFunc, mws ...middleware.Middleware) {
		mux.Handle(pattern, middleware.Chain(fn, append([]middleware.Middleware{global}, mws...)...))
	}

	route("GET /api/admin/billets", h.List, auth, role)
	route("PATCH /api/admin/billets/{id}", h.Update, auth, role)
}
