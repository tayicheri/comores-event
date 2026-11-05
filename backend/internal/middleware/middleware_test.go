package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"comores-event-api/internal/constants"
)

func ok() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
}

func TestCSRF(t *testing.T) {
	h := CSRF()(ok())

	cases := []struct {
		name   string
		method string
		path   string
		cookie string
		header string
		want   int
	}{
		{"GET libre", "GET", "/api/auth/me", "", "", 200},
		{"login exempté", "POST", "/api/auth/login", "", "", 200},
		{"register exempté", "POST", "/api/auth/register", "", "", 200},
		{"refresh exempté", "POST", "/api/auth/refresh", "", "", 200},
		{"logout sans token", "POST", "/api/auth/logout", "", "", 403},
		{"logout token différent", "POST", "/api/auth/logout", "a", "b", 403},
		{"logout ok", "POST", "/api/auth/logout", "abc", "abc", 200},
		{"patch me sans token", "PATCH", "/api/auth/me", "", "", 403},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, nil)
			if c.cookie != "" {
				req.AddCookie(&http.Cookie{Name: constants.CookieCSRF, Value: c.cookie})
			}
			if c.header != "" {
				req.Header.Set(constants.HeaderCSRF, c.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != c.want {
				t.Fatalf("got %d want %d", rec.Code, c.want)
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	h := Recovery()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 500 {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestCorsPreflight(t *testing.T) {
	h := Cors([]string{"http://localhost:3000"})(ok())
	req := httptest.NewRequest("OPTIONS", "/api/auth/me", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "PATCH")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 204 || rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("preflight: %d %v", rec.Code, rec.Header())
	}
}

func TestClientIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.5:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.9, 10.0.0.5")

	if got := ParseTrustedProxies("false").ClientIP(req); got != "10.0.0.5" {
		t.Fatalf("untrusted: %s", got)
	}
	if got := ParseTrustedProxies("10.0.0.0/8").ClientIP(req); got != "203.0.113.9" {
		t.Fatalf("trusted: %s", got)
	}
}
