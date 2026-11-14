package workout

import "context"

// Repository is the outbound port for workout persistence. Implementations live
// in the slice's infrastructure package.
//
// Every query is scoped by profile id, and reads return their exercises already
// attached so the caller never triggers an N+1 query.
type Repository interface {
	// GetByProfile returns every workout of a profile, newest first.
	GetByProfile(ctx context.Context, profileID string) ([]Workout, error)

	// GetByProfileAndDay returns the workouts scheduled for the given lowercase
	// English weekday.
	GetByProfileAndDay(ctx context.Context, profileID, day string) ([]Workout, error)

	// GetByID returns a workout of a profile, with found=false when the id does
	// not belong to it.
	GetByID(ctx context.Context, id, profileID string) (*Workout, bool, error)

	// Save creates a workout and its exercises in one transaction. The workout
	// and exercise ids are expected to be filled in already.
	Save(ctx context.Context, profileID string, w *Workout) error

	// Update replaces a workout and its exercises in one transaction.
	Update(ctx context.Context, id, profileID string, w *Workout) error

	// Delete removes a workout, cascading to its exercises.
	Delete(ctx context.Context, id, profileID string) error
}
