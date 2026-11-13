package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/auth"
)

const seededEmail = "a@b.com"

type fixture struct {
	mux    *http.ServeMux
	users  *mock.UserRepository
	mailer *mock.Mailer
}

// newFixture wires the real router to a service on in-memory mocks, with one
// seeded account (a@b.com).
func newFixture(t *testing.T) *fixture {
	t.Helper()

	users := mock.NewUserRepository()
	users.Seed(&user.User{ID: "u1", Email: seededEmail, Name: "A"}, "old-hash")
	mailer := mock.NewMailer()

	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	svc := service.New(users, mock.NewPasswordResetRepository(), mailer,
		signer, auth.DeriveKey("test-secret", "password-reset"))

	mux := http.NewServeMux()
	Router(mux, signer, svc)
	return &fixture{mux: mux, users: users, mailer: mailer}
}

type response struct {
	status int
	body   string
}

// post sends a JSON POST with no Authorization header.
func (f *fixture) post(path, body string) response {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	f.mux.ServeHTTP(rec, req)
	return response{status: rec.Code, body: strings.TrimSpace(rec.Body.String())}
}

func (f *fixture) forgot(email string) response {
	return f.post("/api/auth/forgot-password", `{"email":"`+email+`"}`)
}

func TestForgotPassword_IdenticalResponses(t *testing.T) {
	f := newFixture(t)

	existing := f.forgot(seededEmail)
	cooldown := f.forgot(seededEmail) // within 60 s of the first
	unknown := f.forgot("nobody@b.com")

	if existing.status != http.StatusOK {
		t.Fatalf("existing email: status = %d, want 200 (no auth required)", existing.status)
	}
	if existing.body == "" {
		t.Fatal("existing email: empty body")
	}
	if unknown != existing {
		t.Errorf("unknown email response = %+v, want %+v", unknown, existing)
	}
	if cooldown != existing {
		t.Errorf("cooldown response = %+v, want %+v", cooldown, existing)
	}
}

func TestForgotPassword_DeliversOnlyForExistingAndOutsideCooldown(t *testing.T) {
	f := newFixture(t)

	f.forgot("nobody@b.com")
	if len(f.mailer.Sent) != 0 {
		t.Fatalf("unknown email delivered %d codes, want 0", len(f.mailer.Sent))
	}

	f.forgot(seededEmail)
	if len(f.mailer.Sent) != 1 {
		t.Fatalf("first request delivered %d codes, want 1", len(f.mailer.Sent))
	}

	f.forgot(seededEmail)
	if len(f.mailer.Sent) != 1 {
		t.Fatalf("cooldown request delivered a code: %d sent, want 1", len(f.mailer.Sent))
	}
}

func TestForgotPassword_NormalizesEmail(t *testing.T) {
	f := newFixture(t)

	res := f.forgot(" A@B.com ")
	if res.status != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.status)
	}
	if len(f.mailer.Sent) != 1 {
		t.Fatalf("delivered %d codes, want 1", len(f.mailer.Sent))
	}
	if got := f.mailer.Sent[0].Email; got != seededEmail {
		t.Errorf("code sent to %q, want %q", got, seededEmail)
	}
}

func TestForgotPassword_MalformedJSON(t *testing.T) {
	f := newFixture(t)

	res := f.post("/api/auth/forgot-password", `{"email":`)
	if res.status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.status)
	}
	if len(f.mailer.Sent) != 0 {
		t.Errorf("malformed request delivered %d codes", len(f.mailer.Sent))
	}
}
