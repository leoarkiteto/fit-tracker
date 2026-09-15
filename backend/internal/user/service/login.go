package service

import (
	"context"
	"log"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/auth"
)

// Login verifies credentials and returns a session.
func (s *Service) Login(ctx context.Context, email, password string) (*user.AuthResponse, error) {
	account, found, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("login: loading account: %v", err)
		return nil, apperror.Internal()
	}
	if !found {
		return nil, user.ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, account.PasswordHash) {
		return nil, user.ErrInvalidCredentials
	}

	// Best effort: a failed stamp must not block the login.
	if err := s.repo.UpdateLastLogin(ctx, account.ID, time.Now()); err != nil {
		log.Printf("login: updating last login: %v", err)
	}

	claims := auth.Claims{UserID: account.ID, Email: account.Email, Name: account.Name}
	if account.ProfileID != nil {
		claims.ProfileID = *account.ProfileID
	}

	token, expiresAt, err := s.signer.GenerateToken(claims)
	if err != nil {
		log.Printf("login: generating token: %v", err)
		return nil, user.ErrGeneratingToken
	}

	return &user.AuthResponse{Token: token, ExpiresAt: expiresAt, User: account.ToDTO()}, nil
}
