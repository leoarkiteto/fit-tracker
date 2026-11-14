package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"fittracker-api/internal/workout"
	"fittracker-api/pkg/apperror"
)

// Create stores a new workout with its exercises. The generated ids are written
// back onto the entity so the handler can echo them in the response.
func (s *Service) Create(ctx context.Context, profileID string, w *workout.Workout) error {
	w.ID = uuid.New().String()
	for i := range w.Exercises {
		w.Exercises[i].ID = uuid.New().String()
	}

	if err := s.repo.Save(ctx, profileID, w); err != nil {
		log.Printf("workouts: creating for profile %s: %v", profileID, err)
		return apperror.Internal()
	}

	return nil
}
