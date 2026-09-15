package service

import (
	"context"
	"log"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/auth"
)

// ChangePassword replaces the password of the authenticated account after
// verifying the current one.
func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if newPassword == "" {
		return user.ErrNewPasswordRequired
	}

	currentHash, found, err := s.repo.GetPasswordHash(ctx, userID)
	if err != nil {
		log.Printf("change password: loading hash: %v", err)
		return apperror.Internal()
	}
	if !found {
		return user.ErrUserNotFound
	}

	if !auth.CheckPasswordHash(currentPassword, currentHash) {
		return user.ErrCurrentPasswordMismatch
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		log.Printf("change password: hashing new password: %v", err)
		return apperror.Internal()
	}

	if err := s.repo.UpdatePassword(ctx, userID, newHash); err != nil {
		log.Printf("change password: updating password: %v", err)
		return apperror.Internal()
	}

	return nil
}
