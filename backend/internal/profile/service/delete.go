package service

import (
	"context"
	"log"

	"fittracker-api/pkg/apperror"
)

// Delete removes a profile and every row that references it. Deleting an
// unknown profile is not an error, matching the previous behaviour.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("profiles: deleting %s: %v", id, err)
		return apperror.Internal()
	}

	return nil
}
