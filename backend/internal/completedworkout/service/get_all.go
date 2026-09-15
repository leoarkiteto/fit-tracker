package service

import (
	"context"
	"log"

	"fittracker-api/internal/completedworkout"
	"fittracker-api/pkg/apperror"
)

// GetAll returns the training history of a profile.
func (s *Service) GetAll(ctx context.Context, profileID string) ([]completedworkout.CompletedWorkout, error) {
	history, err := s.repo.GetByProfile(ctx, profileID)
	if err != nil {
		log.Printf("completed workouts: listing for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	// The rows only carry the ids, so the profile is attached here — this
	// endpoint has always echoed userProfileId on every entry.
	for i := range history {
		history[i].UserProfileID = profileID
	}

	return history, nil
}
