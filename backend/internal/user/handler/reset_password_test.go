package handler

import (
	"net/http"
	"testing"

	"fittracker-api/pkg/auth"
)

const (
	invalidCodeMsg   = "Código inválido ou expirado"
	passwordShortMsg = "A senha deve ter pelo menos 6 caracteres"
)

func (f *fixture) reset(email, code, password string) response {
	return f.post("/api/auth/reset-password",
		`{"email":"`+email+`","code":"`+code+`","newPassword":"`+password+`"}`)
}

// issueCode requests a reset through the router and returns the plaintext code
// the mailer recorded.
func (f *fixture) issueCode(t *testing.T) string {
	t.Helper()
	if res := f.forgot(seededEmail); res.status != http.StatusOK {
		t.Fatalf("forgot-password status = %d, want 200", res.status)
	}
	if len(f.mailer.Sent) == 0 {
		t.Fatal("no code delivered")
	}
	return f.mailer.Sent[len(f.mailer.Sent)-1].Code
}

func TestResetPassword_Success(t *testing.T) {
	f := newFixture(t)
	code := f.issueCode(t)

	res := f.reset(seededEmail, code, "newsecret")
	if res.status != http.StatusNoContent {
		t.Fatalf("status = %d (%q), want 204 (no auth required)", res.status, res.body)
	}

	hash := f.users.Hash("u1")
	if hash == "old-hash" || !auth.CheckPasswordHash("newsecret", hash) {
		t.Error("stored hash was not replaced by the bcrypt hash of the new password")
	}

	// The code is single use.
	again := f.reset(seededEmail, code, "another1")
	if again.status != http.StatusBadRequest || again.body != invalidCodeMsg {
		t.Errorf("reused code response = %+v, want 400 %q", again, invalidCodeMsg)
	}
}

func TestResetPassword_NormalizesEmail(t *testing.T) {
	f := newFixture(t)
	code := f.issueCode(t)

	res := f.reset(" A@B.com ", code, "newsecret")
	if res.status != http.StatusNoContent {
		t.Fatalf("status = %d (%q), want 204", res.status, res.body)
	}
}

func TestResetPassword_WrongCode(t *testing.T) {
	f := newFixture(t)
	code := f.issueCode(t)

	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}
	res := f.reset(seededEmail, wrong, "newsecret")
	if res.status != http.StatusBadRequest || res.body != invalidCodeMsg {
		t.Errorf("response = %+v, want 400 %q", res, invalidCodeMsg)
	}
	if f.users.PasswordUpdates != 0 {
		t.Error("password was updated on a wrong code")
	}
}

func TestResetPassword_UnknownEmail(t *testing.T) {
	f := newFixture(t)

	res := f.reset("nobody@b.com", "123456", "newsecret")
	if res.status != http.StatusBadRequest || res.body != invalidCodeMsg {
		t.Errorf("response = %+v, want 400 %q", res, invalidCodeMsg)
	}
}

func TestResetPassword_NoActiveCode(t *testing.T) {
	f := newFixture(t)

	res := f.reset(seededEmail, "123456", "newsecret")
	if res.status != http.StatusBadRequest || res.body != invalidCodeMsg {
		t.Errorf("response = %+v, want 400 %q", res, invalidCodeMsg)
	}
}

func TestResetPassword_ShortPasswordDoesNotConsumeCode(t *testing.T) {
	f := newFixture(t)
	code := f.issueCode(t)

	res := f.reset(seededEmail, code, "12345")
	if res.status != http.StatusBadRequest || res.body != passwordShortMsg {
		t.Fatalf("response = %+v, want 400 %q", res, passwordShortMsg)
	}

	// The same code still works afterwards: it was neither consumed nor counted.
	ok := f.reset(seededEmail, code, "newsecret")
	if ok.status != http.StatusNoContent {
		t.Errorf("code after short password: status = %d (%q), want 204", ok.status, ok.body)
	}
}

func TestResetPassword_MalformedJSON(t *testing.T) {
	f := newFixture(t)

	res := f.post("/api/auth/reset-password", `not json`)
	if res.status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.status)
	}
}
