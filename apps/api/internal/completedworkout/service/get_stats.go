package service

import (
	"context"
	"log"
	"time"

	"fittracker-api/internal/completedworkout"
	"fittracker-api/pkg/apperror"
)

// statsWindow is how far back "this week" reaches.
const statsWindow = 7 * 24 * time.Hour

// GetStats summarises the training history of a profile. "This week" means the
// last seven days, as before.
func (s *Service) GetStats(ctx context.Context, profileID string) (*completedworkout.Stats, error) {
	since := time.Now().Add(-statsWindow)

	total, err := s.repo.CountByProfile(ctx, profileID)
	if err != nil {
		log.Printf("completed workouts: counting for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	thisWeek, err := s.repo.CountByProfileSince(ctx, profileID, since)
	if err != nil {
		log.Printf("completed workouts: counting the last week for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	seconds, err := s.repo.SumDurationByProfile(ctx, profileID)
	if err != nil {
		log.Printf("completed workouts: summing duration for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	return &completedworkout.Stats{
		TotalWorkoutsCompleted: total,
		WorkoutsThisWeek:       thisWeek,
		TotalMinutesSpent:      seconds / 60,
	}, nil
}
