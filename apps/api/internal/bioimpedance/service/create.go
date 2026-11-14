package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/pkg/apperror"
)

// Create stores a new measurement. A missing date defaults to now, and the
// stored measurement is returned so the handler can echo it with a 201.
func (s *Service) Create(
	ctx context.Context,
	profileID string,
	m *bioimpedance.Measurement,
) (*bioimpedance.Measurement, error) {
	m.ID = uuid.New().String()
	m.UserProfileID = profileID
	if m.Date.IsZero() {
		m.Date = time.Now()
	}

	if err := s.repo.Save(ctx, m); err != nil {
		log.Printf("bioimpedance: storing for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	return m, nil
}
