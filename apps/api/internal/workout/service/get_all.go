package service

import (
	"context"
	"log"

	"fittracker-api/internal/workout"
	"fittracker-api/pkg/apperror"
)

// GetAll returns every workout of a profile.
func (s *Service) GetAll(ctx context.Context, profileID string) ([]workout.Workout, error) {
	workouts, err := s.repo.GetByProfile(ctx, profileID)
	if err != nil {
		log.Printf("workouts: listing for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	// This endpoint has always echoed the profile id back on every workout.
	for i := range workouts {
		workouts[i].UserProfileID = profileID
	}

	return workouts, nil
}
