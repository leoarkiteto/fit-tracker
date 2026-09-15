package service

import (
	"context"
	"log"

	"fittracker-api/internal/profile"
	"fittracker-api/pkg/apperror"
)

// GetByID returns a single profile by id.
func (s *Service) GetByID(ctx context.Context, id string) (*profile.Profile, error) {
	found, ok, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("profiles: loading %s: %v", id, err)
		return nil, apperror.Internal()
	}
	if !ok {
		return nil, profile.ErrProfileNotFound
	}

	return found, nil
}
