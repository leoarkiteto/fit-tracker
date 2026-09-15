package service

import (
	"context"
	"log"

	"fittracker-api/pkg/apperror"
)

// Delete removes a workout belonging to the profile. Deleting an unknown id is
// not an error, matching the previous behaviour.
func (s *Service) Delete(ctx context.Context, id, profileID string) error {
	if err := s.repo.Delete(ctx, id, profileID); err != nil {
		log.Printf("workouts: deleting %s: %v", id, err)
		return apperror.Internal()
	}

	return nil
}
