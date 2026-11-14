package mock

import (
	"context"

	"fittracker-api/internal/aiplanning"
)

// ProfileRepository is an in-memory aiplanning.ProfileRepository.
type ProfileRepository struct {
	GetProfileErr error
	Profiles      map[string]aiplanning.Profile
}

var _ aiplanning.ProfileRepository = (*ProfileRepository)(nil)

// NewProfileRepository builds an empty in-memory repository.
func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{Profiles: map[string]aiplanning.Profile{}}
}

// GetProfile returns the seeded profile.
func (r *ProfileRepository) GetProfile(
	_ context.Context,
	profileID string,
) (*aiplanning.Profile, bool, error) {
	if r.GetProfileErr != nil {
		return nil, false, r.GetProfileErr
	}

	profile, ok := r.Profiles[profileID]
	if !ok {
		return nil, false, nil
	}

	stored := profile
	return &stored, true, nil
}
