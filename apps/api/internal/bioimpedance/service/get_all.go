package service

import (
	"context"
	"log"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/pkg/apperror"
)

// GetAll returns the measurement history of a profile.
func (s *Service) GetAll(ctx context.Context, profileID string) ([]bioimpedance.Measurement, error) {
	history, err := s.repo.GetByProfile(ctx, profileID)
	if err != nil {
		log.Printf("bioimpedance: listing for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	// The rows only carry the id, so the profile is attached here — this
	// endpoint has always echoed userProfileId on every measurement.
	for i := range history {
		history[i].UserProfileID = profileID
	}

	return history, nil
}
