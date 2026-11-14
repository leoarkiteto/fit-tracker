package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/auth"
)

// Register creates a profile + account pair and returns an authenticated session.
func (s *Service) Register(ctx context.Context, email, password, name string) (*user.AuthResponse, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("register: hashing password: %v", err)
		return nil, user.ErrHashingPassword
	}

	profileID := uuid.New().String()
	account := &user.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		ProfileID: &profileID,
	}

	if err := s.repo.Save(ctx, account, hashedPassword); err != nil {
		log.Printf("register: saving account: %v", err)
		if errors.Is(err, user.ErrRegistrationFailed) {
			return nil, user.ErrRegistrationFailed
		}
		return nil, apperror.Internal()
	}

	token, expiresAt, err := s.signer.GenerateToken(auth.Claims{
		UserID:    account.ID,
		Email:     account.Email,
		Name:      account.Name,
		ProfileID: profileID,
	})
	if err != nil {
		log.Printf("register: generating token: %v", err)
		return nil, user.ErrGeneratingToken
	}

	return &user.AuthResponse{Token: token, ExpiresAt: expiresAt, User: account.ToDTO()}, nil
}
