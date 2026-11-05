package auth

import "testing"

func TestJWTRoundTrip(t *testing.T) {
	m := NewJWTManager("0123456789abcdef0123456789abcdef")
	tok, err := m.Generate("user-1", "admin", 3, "sess-1")
	if err != nil {
		t.Fatal(err)
	}
	c, err := m.Parse(tok)
	if err != nil {
		t.Fatal(err)
	}
	if c.UserID() != "user-1" || c.Role != "admin" || c.TokenVersion != 3 || c.SessionID != "sess-1" {
		t.Fatalf("claims inattendus: %+v", c)
	}
	if _, err := NewJWTManager("another-secret-another-secret-123").Parse(tok); err == nil {
		t.Fatal("secret différent accepté")
	}
}

func TestBcryptCompat(t *testing.T) {
	if testing.Short() {
		t.Skip("bcrypt cost 14 lent")
	}
	h, err := HashPassword("secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(h, "secret-pass") || CheckPassword(h, "wrong") {
		t.Fatal("vérification bcrypt incorrecte")
	}
}

func TestOpaqueToken(t *testing.T) {
	a, _ := GenerateOpaqueToken()
	b, _ := GenerateOpaqueToken()
	if a == b || len(HashToken(a)) != 64 {
		t.Fatal("token opaque invalide")
	}
}
