package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/completedworkout"
	"fittracker-api/pkg/apperror"
)

// Complete records a finished workout and returns the stored entry, which the
// handler echoes with a 201.
func (s *Service) Complete(
	ctx context.Context,
	profileID, workoutID string,
	durationSeconds int,
) (*completedworkout.CompletedWorkout, error) {
	entry := &completedworkout.CompletedWorkout{
		ID:              uuid.New().String(),
		WorkoutID:       workoutID,
		UserProfileID:   profileID,
		CompletedAt:     time.Now(),
		DurationSeconds: durationSeconds,
	}

	if err := s.repo.Save(ctx, entry); err != nil {
		log.Printf("completed workouts: recording for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	return entry, nil
}
