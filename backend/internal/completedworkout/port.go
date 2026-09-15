package completedworkout

import (
	"context"
	"time"
)

// Repository is the outbound port for training-history persistence.
// Implementations live in the slice's infrastructure package.
type Repository interface {
	// GetByProfile returns the history of a profile, newest first.
	GetByProfile(ctx context.Context, profileID string) ([]CompletedWorkout, error)

	// Save records a completed workout. Ids and timestamps are expected to be
	// filled in already.
	Save(ctx context.Context, c *CompletedWorkout) error

	// CountByProfile counts every completed workout of a profile.
	CountByProfile(ctx context.Context, profileID string) (int, error)

	// CountByProfileSince counts the workouts completed at or after an instant.
	CountByProfileSince(ctx context.Context, profileID string, since time.Time) (int, error)

	// SumDurationByProfile totals the recorded duration in seconds, returning 0
	// for a profile with no history.
	SumDurationByProfile(ctx context.Context, profileID string) (int, error)
}
