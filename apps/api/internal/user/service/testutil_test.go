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

// socialFixture bundles the doubles wired by newSocialTestService.
type socialFixture struct {
	svc        *Service
	signer     *auth.Signer
	users      *mock.UserRepository
	resets     *mock.PasswordResetRepository
	mailer     *mock.Mailer
	identities *mock.IdentityRepository
	google     *mock.GoogleVerifier
	github     *mock.GitHubExchanger
}

// newSocialTestService builds the service with social login enabled over
// in-memory doubles that share state.
func newSocialTestService(t *testing.T) *socialFixture {
	t.Helper()

	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("building signer: %v", err)
	}

	users := mock.NewUserRepository()
	resets := mock.NewPasswordResetRepository()
	mailer := mock.NewMailer()
	identities := mock.NewIdentityRepository(users, resets)
	google := &mock.GoogleVerifier{}
	github := &mock.GitHubExchanger{}

	svc := New(users, resets, mailer, signer, auth.DeriveKey("test-secret", "password-reset")).
		WithSocialLogin(identities, google, github)

	return &socialFixture{svc, signer, users, resets, mailer, identities, google, github}
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
