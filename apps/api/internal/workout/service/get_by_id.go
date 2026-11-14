package service

import (
	"context"
	"log"

	"fittracker-api/internal/workout"
	"fittracker-api/pkg/apperror"
)

// GetByID returns a single workout belonging to the profile.
func (s *Service) GetByID(ctx context.Context, id, profileID string) (*workout.Workout, error) {
	found, ok, err := s.repo.GetByID(ctx, id, profileID)
	if err != nil {
		log.Printf("workouts: loading %s: %v", id, err)
		return nil, apperror.Internal()
	}
	if !ok {
		return nil, workout.ErrWorkoutNotFound
	}

	found.UserProfileID = profileID
	return found, nil
}
