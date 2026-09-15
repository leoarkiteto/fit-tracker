package service

import (
	"context"
	"log"

	"fittracker-api/internal/profile"
	"fittracker-api/pkg/apperror"
)

// Update stores the editable fields of a profile.
func (s *Service) Update(ctx context.Context, p *profile.Profile) error {
	if err := s.repo.Update(ctx, p); err != nil {
		log.Printf("profiles: updating %s: %v", p.ID, err)
		return apperror.Internal()
	}

	return nil
}
