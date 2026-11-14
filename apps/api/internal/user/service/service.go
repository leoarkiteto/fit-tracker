// Package service holds the account use cases. It depends only on the domain
// and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/auth"
)

// Service holds the account use cases.
type Service struct {
	repo     user.Repository
	resets   user.PasswordResetRepository
	mailer   user.Mailer
	signer   *auth.Signer
	now      func() time.Time
	resetKey []byte

	// Social login ports; nil until WithSocialLogin is called.
	identities user.IdentityRepository
	google     user.GoogleVerifier
	github     user.GitHubExchanger
}

// WithSocialLogin enables social login. Any nil argument leaves the matching
// provider unavailable.
func (s *Service) WithSocialLogin(identities user.IdentityRepository, google user.GoogleVerifier,
	github user.GitHubExchanger) *Service {
	s.identities = identities
	s.google = google
	s.github = github
	return s
}

// New builds the account service. resetKey is the HMAC key under which reset
// codes are hashed before they are stored.
func New(repo user.Repository, resets user.PasswordResetRepository, mailer user.Mailer,
	signer *auth.Signer, resetKey []byte) *Service {
	return &Service{
		repo:     repo,
		resets:   resets,
		mailer:   mailer,
		signer:   signer,
		now:      time.Now,
		resetKey: resetKey,
	}
}
