package service

import (
	"context"
	"log"

	"fittracker-api/internal/profile"
	"fittracker-api/pkg/apperror"
)

// GetAll returns every profile.
func (s *Service) GetAll(ctx context.Context) ([]profile.Profile, error) {
	profiles, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("profiles: listing: %v", err)
		return nil, apperror.Internal()
	}

	return profiles, nil
}
