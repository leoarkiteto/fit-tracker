package profile

import "context"

// Repository is the outbound port for profile persistence. Implementations live
// in the slice's infrastructure package.
type Repository interface {
	// GetAll returns every profile. An empty result is a nil slice, which the
	// API has always serialised as JSON null.
	GetAll(ctx context.Context) ([]Profile, error)

	// GetByID returns a profile, with found=false when it does not exist.
	GetByID(ctx context.Context, id string) (*Profile, bool, error)

	// Update stores every editable field of a profile.
	Update(ctx context.Context, p *Profile) error

	// Delete removes a profile together with every row referencing it.
	Delete(ctx context.Context, id string) error
}
