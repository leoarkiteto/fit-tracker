package service

import (
	"context"
	"log"
	"strings"
	"time"

	"fittracker-api/internal/workout"
	"fittracker-api/pkg/apperror"
)

// GetToday returns the workouts scheduled for the current weekday.
func (s *Service) GetToday(ctx context.Context, profileID string) ([]workout.Workout, error) {
	// Days are stored comma-separated and matched against the lowercase English
	// weekday name.
	today := strings.ToLower(time.Now().Weekday().String())

	workouts, err := s.repo.GetByProfileAndDay(ctx, profileID, today)
	if err != nil {
		log.Printf("workouts: listing today's for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	// Unlike GetAll and GetByID, this endpoint has never echoed the profile id.
	return workouts, nil
}
