package service

import (
	"context"
	"log"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
)

// Me returns the public projection of the authenticated account.
func (s *Service) Me(ctx context.Context, userID string) (*user.DTO, error) {
	account, found, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("me: loading account: %v", err)
		return nil, apperror.Internal()
	}
	if !found {
		return nil, user.ErrUserNotFound
	}

	dto := account.ToDTO()
	return &dto, nil
}
