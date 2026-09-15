package service

import (
	"context"
	"log"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/pkg/apperror"
)

// GetLatest returns the most recent measurement of a profile.
func (s *Service) GetLatest(ctx context.Context, profileID string) (*bioimpedance.Measurement, error) {
	found, ok, err := s.repo.GetLatest(ctx, profileID)
	if err != nil {
		log.Printf("bioimpedance: loading latest for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}
	if !ok {
		return nil, bioimpedance.ErrNoMeasurement
	}

	found.UserProfileID = profileID
	return found, nil
}
