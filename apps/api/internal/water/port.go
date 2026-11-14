package water

import (
	"context"
	"time"
)

// Repository is the outbound port for hydration persistence. Implementations
// live in the slice's infrastructure package.
type Repository interface {
	// GetByDay returns the entries of a profile logged in [start, end), in the
	// order they were recorded.
	GetByDay(ctx context.Context, profileID string, start, end time.Time) ([]Entry, error)

	// Save stores an entry. Its id is expected to be filled in.
	Save(ctx context.Context, e *Entry) error

	// Delete removes an entry belonging to a profile. Deleting an unknown id is
	// not an error.
	Delete(ctx context.Context, id, profileID string) error

	// CurrentWeightKg returns the recorded weight of a profile, with
	// found=false when it has no weight or does not exist.
	CurrentWeightKg(ctx context.Context, profileID string) (float64, bool, error)
}
