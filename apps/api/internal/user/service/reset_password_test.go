package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/pkg/auth"
)

// issueCode seeds an account, requests a reset at t0 and returns the delivered
// code together with a clock the test can advance.
func issueCode(t *testing.T) (*Service, *mock.UserRepository, *mock.PasswordResetRepository, *mock.Mailer, string, *time.Time) {
	t.Helper()
	svc, repo, resets, mailer := newResetTestService(t)
	now := t0
	svc.now = func() time.Time { return now }
	seedAccount(repo, "a@b.com")
	if err := svc.RequestPasswordReset(context.Background(), "a@b.com"); err != nil {
		t.Fatal(err)
	}
	return svc, repo, resets, mailer, mailer.Sent[0].Code, &now
}

// wrongCode returns a 6-digit code different from code.
func wrongCode(code string) string {
	if code == "000000" {
		return "000001"
	}
	return "000000"
}

func requireInvalid(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, user.ErrInvalidResetCode) {
		t.Fatalf("expected ErrInvalidResetCode, got %v", err)
	}
}

func TestResetPassword_SuccessReplacesPasswordAndCodeIsSingleUse(t *testing.T) {
	svc, repo, resets, _, code, _ := issueCode(t)
	ctx := context.Background()

	if err := svc.ResetPassword(ctx, "a@b.com", code, "newsecret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hash := repo.Hash("u1")
	if hash == "old-hash" || !auth.CheckPasswordHash("newsecret", hash) {
		t.Fatal("password was not replaced by a bcrypt hash of the new one")
	}
	if stored, _ := resets.Code("u1"); stored.UsedAt == nil {
		t.Error("code not marked used")
	}

	requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", code, "another1"))
	if repo.Hash("u1") != hash {
		t.Error("reused code changed the password")
	}
}

func TestResetPassword_NewCodeInvalidatesOldOne(t *testing.T) {
	svc, repo, _, mailer, codeA, now := issueCode(t)
	ctx := context.Background()

	*now = t0.Add(2 * user.ResetCodeCooldown)
	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatal(err)
	}
	codeB := mailer.Sent[1].Code
	if codeA == codeB {
		t.Skip("random codes collided")
	}

	requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", codeA, "newsecret"))
	if err := svc.ResetPassword(ctx, "a@b.com", codeB, "newsecret"); err != nil {
		t.Fatalf("code B should succeed: %v", err)
	}
	if repo.Hash("u1") == "old-hash" {
		t.Error("password unchanged")
	}
}

func TestResetPassword_CodeFromCooldownRequestStillWorks(t *testing.T) {
	svc, _, _, _, code, now := issueCode(t)
	ctx := context.Background()

	*now = t0.Add(30 * time.Second)
	if err := svc.RequestPasswordReset(ctx, "a@b.com"); err != nil {
		t.Fatal(err)
	}
	if err := svc.ResetPassword(ctx, "a@b.com", code, "newsecret"); err != nil {
		t.Fatalf("original code should still work: %v", err)
	}
}

func TestResetPassword_ExpiredCode(t *testing.T) {
	svc, repo, resets, _, code, now := issueCode(t)

	*now = t0.Add(16 * time.Minute)
	requireInvalid(t, svc.ResetPassword(context.Background(), "a@b.com", code, "newsecret"))
	if repo.Hash("u1") != "old-hash" {
		t.Error("password changed")
	}
	if stored, _ := resets.Code("u1"); stored.Attempts != 0 {
		t.Error("expired code should not count an attempt")
	}
}

func TestResetPassword_TTLBoundaryIsExclusive(t *testing.T) {
	svc, _, _, _, code, now := issueCode(t)

	*now = t0.Add(user.ResetCodeTTL)
	requireInvalid(t, svc.ResetPassword(context.Background(), "a@b.com", code, "newsecret"))
}

func TestResetPassword_FourWrongThenCorrectSucceeds(t *testing.T) {
	svc, repo, _, _, code, _ := issueCode(t)
	ctx := context.Background()

	for i := 0; i < user.ResetCodeMaxAttempts-1; i++ {
		requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", wrongCode(code), "newsecret"))
	}
	if err := svc.ResetPassword(ctx, "a@b.com", code, "newsecret"); err != nil {
		t.Fatalf("correct code after 4 wrong should succeed: %v", err)
	}
	if repo.Hash("u1") == "old-hash" {
		t.Error("password unchanged")
	}
}

func TestResetPassword_FiveWrongThenCorrectFails(t *testing.T) {
	svc, repo, resets, _, code, _ := issueCode(t)
	ctx := context.Background()

	for i := 0; i < user.ResetCodeMaxAttempts; i++ {
		requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", wrongCode(code), "newsecret"))
	}
	requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", code, "newsecret"))
	if repo.Hash("u1") != "old-hash" {
		t.Error("password changed")
	}
	if stored, _ := resets.Code("u1"); stored.Attempts != user.ResetCodeMaxAttempts {
		t.Errorf("attempts = %d", stored.Attempts)
	}
}

func TestResetPassword_UnknownEmailAndNoActiveCode(t *testing.T) {
	svc, repo, _, _ := newResetTestService(t)
	seedAccount(repo, "a@b.com")
	ctx := context.Background()

	requireInvalid(t, svc.ResetPassword(ctx, "nobody@example.com", "123456", "newsecret"))
	requireInvalid(t, svc.ResetPassword(ctx, "a@b.com", "123456", "newsecret"))
}

func TestResetPassword_ShortPasswordDoesNotConsumeAttempt(t *testing.T) {
	svc, repo, resets, _, code, _ := issueCode(t)
	ctx := context.Background()

	err := svc.ResetPassword(ctx, "a@b.com", code, "12345")
	if !errors.Is(err, user.ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
	if stored, _ := resets.Code("u1"); stored.Attempts != 0 || stored.UsedAt != nil {
		t.Errorf("short password touched the code: %+v", stored)
	}
	if err := svc.ResetPassword(ctx, "a@b.com", code, "newsecret"); err != nil {
		t.Fatalf("same code should then succeed: %v", err)
	}
	if repo.Hash("u1") == "old-hash" {
		t.Error("password unchanged")
	}
}

func TestResetPassword_NormalizesEmail(t *testing.T) {
	svc, repo, _, _, code, _ := issueCode(t)

	if err := svc.ResetPassword(context.Background(), " A@B.com ", code, "newsecret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Hash("u1") == "old-hash" {
		t.Error("password unchanged")
	}
}

func TestResetPassword_StoredHashIsDerivedHMAC(t *testing.T) {
	svc, _, resets, _, code, _ := issueCode(t)

	stored, _ := resets.Code("u1")
	if stored.CodeHash == code {
		t.Fatal("plaintext stored")
	}
	if stored.CodeHash != svc.hashResetCode(code) {
		t.Fatal("stored hash is not HMAC(derived key, code)")
	}
}

func TestResetPassword_RepositoryErrorsAreInternal(t *testing.T) {
	boom := errors.New("boom")
	tests := map[string]func(*mock.UserRepository, *mock.PasswordResetRepository){
		"get by email":    func(r *mock.UserRepository, _ *mock.PasswordResetRepository) { r.GetByEmailErr = boom },
		"get code":        func(_ *mock.UserRepository, c *mock.PasswordResetRepository) { c.GetErr = boom },
		"consume attempt": func(_ *mock.UserRepository, c *mock.PasswordResetRepository) { c.ConsumeAttemptErr = boom },
		"mark used":       func(_ *mock.UserRepository, c *mock.PasswordResetRepository) { c.MarkUsedErr = boom },
		"update password": func(r *mock.UserRepository, _ *mock.PasswordResetRepository) { r.UpdatePasswordErr = boom },
	}
	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			svc, repo, resets, _, code, _ := issueCode(t)
			setup(repo, resets)

			err := svc.ResetPassword(context.Background(), "a@b.com", code, "newsecret")
			if err == nil || err.Error() != "Internal server error" {
				t.Fatalf("expected internal error, got %v", err)
			}
		})
	}
}
