package service

import (
	"context"
	"crypto/hmac"
	"log"
	"strings"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/auth"
)

// ResetPassword replaces the account password when code is the account's valid
// reset code. Every failure caused by the code or the account yields
// user.ErrInvalidResetCode. ConsumeAttempt is the only attempt gate.
func (s *Service) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	if len(newPassword) < user.MinPasswordLength {
		return user.ErrPasswordTooShort
	}

	account, found, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("reset password: loading account: %v", err)
		return apperror.Internal()
	}
	if !found {
		return user.ErrInvalidResetCode
	}

	now := s.now().UTC()

	stored, exists, err := s.resets.Get(ctx, account.ID)
	if err != nil {
		log.Printf("reset password: loading code: %v", err)
		return apperror.Internal()
	}
	if !exists || stored.UsedAt != nil || now.Sub(stored.CreatedAt) >= user.ResetCodeTTL {
		return user.ErrInvalidResetCode
	}

	hash := s.hashResetCode(code)

	granted, err := s.resets.ConsumeAttempt(ctx, account.ID, stored.CodeHash, user.ResetCodeMaxAttempts)
	if err != nil {
		log.Printf("reset password: consuming attempt: %v", err)
		return apperror.Internal()
	}
	if !granted {
		return user.ErrInvalidResetCode
	}

	if !hmac.Equal([]byte(hash), []byte(stored.CodeHash)) {
		return user.ErrInvalidResetCode
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		log.Printf("reset password: hashing password: %v", err)
		return apperror.Internal()
	}

	won, err := s.resets.MarkUsed(ctx, account.ID, stored.CodeHash, now)
	if err != nil {
		log.Printf("reset password: marking code used: %v", err)
		return apperror.Internal()
	}
	if !won {
		return user.ErrInvalidResetCode
	}

	if err := s.repo.UpdatePassword(ctx, account.ID, newHash); err != nil {
		log.Printf("reset password: updating password: %v", err)
		return apperror.Internal()
	}

	return nil
}
