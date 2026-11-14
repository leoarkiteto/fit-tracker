package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/auth"
)

// LoginWithGoogle verifies a Google ID token and signs the account in.
func (s *Service) LoginWithGoogle(ctx context.Context, idToken string) (*user.AuthResponse, error) {
	if s.google == nil || s.identities == nil {
		return nil, user.ErrSocialUnavailable(user.ProviderGoogle)
	}

	identity, err := s.google.Verify(ctx, idToken)
	if err != nil || identity == nil {
		// The error and the token are never logged: they may carry provider payloads.
		log.Printf("social login: google verification failed")
		return nil, user.ErrSocialLoginFailed(user.ProviderGoogle)
	}

	return s.socialLogin(ctx, user.ProviderGoogle, identity)
}

// LoginWithGitHub exchanges a GitHub authorization code and signs the account in.
func (s *Service) LoginWithGitHub(ctx context.Context, code, codeVerifier, redirectURI string) (*user.AuthResponse, error) {
	if s.github == nil || s.identities == nil {
		return nil, user.ErrSocialUnavailable(user.ProviderGitHub)
	}

	identity, err := s.github.Exchange(ctx, code, codeVerifier, redirectURI)
	if err != nil || identity == nil {
		log.Printf("social login: github exchange failed")
		return nil, user.ErrSocialLoginFailed(user.ProviderGitHub)
	}

	return s.socialLogin(ctx, user.ProviderGitHub, identity)
}

// socialLogin resolves the account for a verified provider identity, creating
// or linking it when needed, and issues a session.
func (s *Service) socialLogin(ctx context.Context, provider string, identity *user.ProviderIdentity) (*user.AuthResponse, error) {
	account, found, err := s.identities.FindUserByIdentity(ctx, provider, identity.Subject)
	if err != nil {
		log.Printf("social login: finding identity (%s): %v", provider, err)
		return nil, user.ErrSocialLoginFailed(provider)
	}

	if !found {
		account, err = s.createOrLink(ctx, provider, identity)
		if err != nil {
			return nil, err
		}
	}

	// Best effort: a failed stamp must not block the login.
	if err := s.repo.UpdateLastLogin(ctx, account.ID, time.Now()); err != nil {
		log.Printf("social login: updating last login: %v", err)
	}

	claims := auth.Claims{UserID: account.ID, Email: account.Email, Name: account.Name}
	if account.ProfileID != nil {
		claims.ProfileID = *account.ProfileID
	}

	token, expiresAt, err := s.signer.GenerateToken(claims)
	if err != nil {
		log.Printf("social login: generating token: %v", err)
		return nil, user.ErrGeneratingToken
	}

	return &user.AuthResponse{Token: token, ExpiresAt: expiresAt, User: account.ToDTO()}, nil
}

// createOrLink handles an identity that is not linked yet.
func (s *Service) createOrLink(ctx context.Context, provider string, identity *user.ProviderIdentity) (*user.User, error) {
	email := strings.ToLower(strings.TrimSpace(identity.Email))
	if !identity.EmailVerified || email == "" {
		return nil, user.ErrSocialNoVerifiedEmail(provider)
	}

	matches, err := s.identities.FindUsersByEmailFold(ctx, email)
	if err != nil {
		log.Printf("social login: finding accounts by email (%s): %v", provider, err)
		return nil, user.ErrSocialLoginFailed(provider)
	}

	linked := user.LinkedIdentity{
		Provider:  provider,
		Subject:   identity.Subject,
		CreatedAt: s.now().UTC(),
	}

	switch len(matches) {
	case 0:
		profileID := uuid.New().String()
		account := &user.User{
			ID:        uuid.New().String(),
			ProfileID: &profileID,
			Email:     email,
			Name:      socialDisplayName(identity, email),
		}
		linked.UserID = account.ID
		if err := s.identities.CreateWithIdentity(ctx, account, linked); err != nil {
			log.Printf("social login: creating account (%s): %v", provider, err)
			return nil, user.ErrSocialLoginFailed(provider)
		}
		return account, nil

	case 1:
		account := matches[0]
		linked.UserID = account.ID
		if err := s.identities.LinkIdentity(ctx, linked); err != nil {
			log.Printf("social login: linking identity (%s): %v", provider, err)
			return nil, user.ErrSocialLoginFailed(provider)
		}
		// What the repository stored: lowercased email.
		account.Email = strings.ToLower(account.Email)
		return &account, nil

	default:
		log.Printf("social login: %d accounts share the verified email (%s)", len(matches), provider)
		return nil, user.ErrSocialLoginFailed(provider)
	}
}

// socialDisplayName picks the name of a new account: provider name, then
// login, then the local part of the email.
func socialDisplayName(identity *user.ProviderIdentity, email string) string {
	if name := strings.TrimSpace(identity.Name); name != "" {
		return name
	}
	if login := strings.TrimSpace(identity.Login); login != "" {
		return login
	}
	local, _, _ := strings.Cut(email, "@")
	return strings.TrimSpace(local)
}
