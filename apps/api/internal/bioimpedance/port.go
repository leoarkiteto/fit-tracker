package bioimpedance

import "context"

// Repository is the outbound port for bioimpedance persistence. Implementations
// live in the slice's infrastructure package.
type Repository interface {
	// GetByProfile returns every measurement of a profile, newest first.
	GetByProfile(ctx context.Context, profileID string) ([]Measurement, error)

	// GetLatest returns the most recent measurement, with found=false when the
	// profile has none.
	GetLatest(ctx context.Context, profileID string) (*Measurement, bool, error)

	// Save stores a measurement. Its id and date are expected to be filled in.
	Save(ctx context.Context, m *Measurement) error

	// Delete removes a measurement belonging to a profile. Deleting an unknown
	// id is not an error.
	Delete(ctx context.Context, id, profileID string) error
}
