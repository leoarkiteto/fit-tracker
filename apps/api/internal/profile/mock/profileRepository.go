// Package mock holds hand-written test doubles for the profile slice, filling
// the role mockery-generated mocks play in the reference architecture.
package mock

import (
	"context"

	"fittracker-api/internal/profile"
)

// ProfileRepository is an in-memory profile.Repository. Set any Err field to
// force that call to fail; the counters record what the service did.
type ProfileRepository struct {
	GetAllErr  error
	GetByIDErr error
	UpdateErr  error
	DeleteErr  error

	UpdateCalls int
	DeleteCalls int

	profiles map[string]*profile.Profile
}

var _ profile.Repository = (*ProfileRepository)(nil)

// NewProfileRepository builds an empty in-memory repository.
func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{profiles: map[string]*profile.Profile{}}
}

// Seed registers a profile that appears to already exist.
func (r *ProfileRepository) Seed(p *profile.Profile) {
	stored := *p
	r.profiles[p.ID] = &stored
}

// GetAll returns every stored profile, or a nil slice when there is none — the
// same shape the SQL implementation returns for an empty table.
func (r *ProfileRepository) GetAll(_ context.Context) ([]profile.Profile, error) {
	if r.GetAllErr != nil {
		return nil, r.GetAllErr
	}

	if len(r.profiles) == 0 {
		return nil, nil
	}

	profiles := make([]profile.Profile, 0, len(r.profiles))
	for _, p := range r.profiles {
		profiles = append(profiles, *p)
	}

	return profiles, nil
}

// GetByID returns the profile with the given id.
func (r *ProfileRepository) GetByID(_ context.Context, id string) (*profile.Profile, bool, error) {
	if r.GetByIDErr != nil {
		return nil, false, r.GetByIDErr
	}

	p, ok := r.profiles[id]
	if !ok {
		return nil, false, nil
	}

	stored := *p
	return &stored, true, nil
}

// Update stores the profile fields.
func (r *ProfileRepository) Update(_ context.Context, p *profile.Profile) error {
	if r.UpdateErr != nil {
		return r.UpdateErr
	}

	stored := *p
	r.profiles[p.ID] = &stored
	r.UpdateCalls++
	return nil
}

// Delete removes a profile.
func (r *ProfileRepository) Delete(_ context.Context, id string) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}

	delete(r.profiles, id)
	r.DeleteCalls++
	return nil
}
