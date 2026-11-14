package service

import (
	"context"
	"log"
	"time"

	"fittracker-api/internal/water"
	"fittracker-api/pkg/apperror"
)

// dateLayout is the format the ?date= parameter uses.
const dateLayout = "2006-01-02"

// GetDaily returns the hydration summary of a profile for one day. A missing
// date means today; days are bounded in UTC.
func (s *Service) GetDaily(ctx context.Context, profileID, dateParam string) (*water.DailySummary, error) {
	date := time.Now()
	if dateParam != "" {
		parsed, err := time.Parse(dateLayout, dateParam)
		if err != nil {
			return nil, water.ErrInvalidDate
		}
		date = parsed
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	entries, err := s.repo.GetByDay(ctx, profileID, start, end)
	if err != nil {
		log.Printf("water: listing for profile %s: %v", profileID, err)
		return nil, apperror.Internal()
	}

	totalMl := 0
	for _, entry := range entries {
		totalMl += entry.AmountMl
	}

	return &water.DailySummary{
		Date:    start.Format(dateLayout),
		TotalMl: totalMl,
		GoalMl:  s.dailyGoal(ctx, profileID),
		Entries: entries,
	}, nil
}

// dailyGoal derives the daily target from the profile weight. A failing lookup
// falls back to the default weight instead of failing the request, which is the
// behaviour the previous implementation had.
func (s *Service) dailyGoal(ctx context.Context, profileID string) int {
	weight, found, err := s.repo.CurrentWeightKg(ctx, profileID)
	if err != nil {
		log.Printf("water: loading weight for profile %s: %v", profileID, err)
		found = false
	}
	if !found {
		weight = water.DefaultWeightKg
	}

	return int(weight * water.MlPerKg)
}
