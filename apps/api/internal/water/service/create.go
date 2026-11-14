package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/water"
	"fittracker-api/pkg/apperror"
)

// Create logs a drink. The stored entry is returned so the handler can echo it
// with a 201.
func (s *Service) Create(
	ctx context.Context,
	profileID string,
	amountMl int,
	consumedAt *time.Time,
) (*water.Entry, error) {
	entry := &water.Entry{
		ID:            uuid.New().String(),
		UserProfileID: profileID,
		AmountMl:      amountMl,
		ConsumedAt:    time.Now(),
	}
	if consumedAt != nil {
		entry.ConsumedAt = *consumedAt
	}

	if err := s.repo.Save(ctx, entry); err != nil {
		log.Printf("water: storing entry for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	return entry, nil
}
