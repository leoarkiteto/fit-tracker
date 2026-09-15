package service

import (
	"context"
	"log"

	"fittracker-api/internal/aiplanning"
	"fittracker-api/pkg/apperror"
)

// Accept stores the workouts of a plan the user accepted for their profile.
func (s *Service) Accept(
	ctx context.Context,
	profileID string,
	planned []aiplanning.PlannedWorkout,
) error {
	if err := s.writer.CreateAll(ctx, profileID, planned); err != nil {
		log.Printf("ai planning: storing accepted plan for profile %s: %v", profileID, err)
		return apperror.Internal()
	}

	return nil
}
