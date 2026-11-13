package service

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
)

var t0 = time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)

// seedAccount registers an account with a placeholder password hash.
func seedAccount(repo *mock.UserRepository, email string) *user.User {
	u := &user.User{ID: "u1", Email: email, Name: "Ana"}
	repo.Seed(u, "old-hash")
	return u
}

func TestRequestPasswordReset_UnknownEmailDeliversNothing(t *testing.T) {
	svc, _, resets, mailer := newResetTestService(t)

	if err := svc.RequestPasswordReset(context.Background(), "nobody@example.com"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(mailer.Sent) != 0 || resets.Replaces != 0 {
		t.Fatalf("expected nothing sent or stored, got %d sent, %d stored", len(mailer.Sent), resets.Replaces)
	}
}

func TestRequestPasswordReset_DeliversSixDigitCodeAndStoresOnlyHash(t *testing.T) {
	svc, repo, resets, mailer := newResetTestService(t)
	svc.now = func() time.Time { return t0 }
	seedAccount(repo, "a@b.com")

	if err := svc.RequestPasswordReset(context.Background(), "a@b.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mailer.Sent) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(mailer.Sent))
	}
	sent := mailer.Sent[0]
	if sent.Email != "a@b.com" {
		t.Errorf("sent to %q", sent.Email)
	}
	if !regexp.MustCompile(`^\d{6}$`).MatchString(sent.Code) {
		t.Errorf("code %q is not 6 digits", sent.Code)
	}

	stored, ok := resets.Code("u1")
	if !ok {
		t.Fatal("no code stored")
	}
	if stored.CodeHash == sent.Code {
		t.Error("plaintext code was stored")
	}
	if stored.CodeHash != svc.hashResetCode(sent.Code) {
		t.Error("stored value is not the HMAC of the code under the derived key")
	}
	if len(stored.CodeHash) != 64 {
		t.Errorf("expected hex sha256, got len %d", len(stored.CodeHash))
	}
	if !stored.CreatedAt.Equal(t0) || stored.Attempts != 0 || stored.UsedAt != nil {
		t.Errorf("unexpected stored code: %+v", stored)
	}
}

func TestRequestPasswordReset_NormalizesEmail(t *testing.T) {
	svc, repo, _, mailer := newResetTestService(t)
	seedAccount(repo, "a@b.com")

	if err := svc.RequestPasswordReset(context.Background(), " A@B.com "); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mailer.Sent) != 1 || mailer.Sent[0].Email != "a@b.com" {
		t.Fatalf("expected delivery to a@b.com, got %+v", mailer.Sent)
	}
}

func TestRequestPasswordReset_CooldownSkipsAndKeepsPreviousCode(t *testing.T) {
	svc, repo, resets, mailer := newResetTestService(t)
	now := t0
	svc.now = func() time.Time { return now }
	seedAccount(repo, "a@b.com")
	ctx := context.Background()

	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatal(err)
	}
	first := mailer.Sent[0].Code

	now = t0.Add(30 * time.Second)
	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatalf("expected nil within cooldown, got %v", err)
	}
	if len(mailer.Sent) != 1 || resets.Replaces != 1 {
		t.Fatalf("expected no new code, got %d sent, %d replaces", len(mailer.Sent), resets.Replaces)
	}
	stored, _ := resets.Code("u1")
	if stored.CodeHash != svc.hashResetCode(first) {
		t.Error("previous code was replaced")
	}
}

func TestRequestPasswordReset_NewCodeAfterCooldownInvalidatesOld(t *testing.T) {
	svc, repo, resets, mailer := newResetTestService(t)
	now := t0
	svc.now = func() time.Time { return now }
	seedAccount(repo, "a@b.com")
	ctx := context.Background()

	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatal(err)
	}
	now = t0.Add(user.ResetCodeCooldown)
	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatal(err)
	}
	if len(mailer.Sent) != 2 || resets.Replaces != 2 {
		t.Fatalf("expected a second code, got %d sent", len(mailer.Sent))
	}
	stored, _ := resets.Code("u1")
	if stored.CodeHash != svc.hashResetCode(mailer.Sent[1].Code) {
		t.Error("stored code is not the newest")
	}
	if !stored.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", stored.CreatedAt, now)
	}
}

func TestRequestPasswordReset_MailerFailureIsLoggedNotReturned(t *testing.T) {
	svc, repo, resets, mailer := newResetTestService(t)
	seedAccount(repo, "a@b.com")
	mailer.SendErr = errors.New("smtp down")

	if err := svc.RequestPasswordReset(context.Background(), "a@b.com"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if _, ok := resets.Code("u1"); !ok {
		t.Error("code should still be stored")
	}
}

func TestRequestPasswordReset_RepositoryErrorsAreInternal(t *testing.T) {
	boom := errors.New("boom")
	tests := map[string]func(*mock.UserRepository, *mock.PasswordResetRepository){
		"get by email": func(r *mock.UserRepository, _ *mock.PasswordResetRepository) { r.GetByEmailErr = boom },
		"get code":     func(_ *mock.UserRepository, c *mock.PasswordResetRepository) { c.GetErr = boom },
		"replace":      func(_ *mock.UserRepository, c *mock.PasswordResetRepository) { c.ReplaceErr = boom },
	}
	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			svc, repo, resets, mailer := newResetTestService(t)
			seedAccount(repo, "a@b.com")
			setup(repo, resets)

			err := svc.RequestPasswordReset(context.Background(), "a@b.com")
			if err == nil || err.Error() != "Internal server error" {
				t.Fatalf("expected internal error, got %v", err)
			}
			if len(mailer.Sent) != 0 {
				t.Error("nothing should be delivered")
			}
		})
	}
}
