package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"fittracker-api/internal/workout"
	"fittracker-api/pkg/apperror"
)

// Update replaces a workout and its exercises. As before, the exercises are
// re-created with fresh ids and the response echoes the submitted body.
func (s *Service) Update(ctx context.Context, id, profileID string, w *workout.Workout) error {
	w.ID = id
	for i := range w.Exercises {
		w.Exercises[i].ID = uuid.New().String()
	}

	if err := s.repo.Update(ctx, id, profileID, w); err != nil {
		log.Printf("workouts: updating %s: %v", id, err)
		return apperror.Internal()
	}

	return nil
}
