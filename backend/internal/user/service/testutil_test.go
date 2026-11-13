package service

import (
	"testing"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/pkg/auth"
)

// newTestService builds the service under test with an in-memory repository and
// a signer that does not depend on the JWT_SECRET environment variable. The
// reset repository and mailer are empty mocks.
func newTestService(t *testing.T, repo user.Repository) (*Service, *auth.Signer) {
	t.Helper()

	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("building signer: %v", err)
	}

	svc := New(repo, mock.NewPasswordResetRepository(), mock.NewMailer(), signer,
		auth.DeriveKey("test-secret", "password-reset"))
	return svc, signer
}

// newResetTestService builds the service with in-memory doubles for every port
// and returns them so password reset tests can seed and inspect them.
func newResetTestService(t *testing.T) (*Service, *mock.UserRepository, *mock.PasswordResetRepository, *mock.Mailer) {
	t.Helper()

	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("building signer: %v", err)
	}

	repo := mock.NewUserRepository()
	resets := mock.NewPasswordResetRepository()
	mailer := mock.NewMailer()
	svc := New(repo, resets, mailer, signer, auth.DeriveKey("test-secret", "password-reset"))
	return svc, repo, resets, mailer
}
