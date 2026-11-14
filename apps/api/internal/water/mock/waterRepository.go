// Package mock holds hand-written test doubles for the hydration slice.
package mock

import (
	"context"
	"time"

	"fittracker-api/internal/water"
)

// EntryRepository is an in-memory water.Repository. Set any Err field to force
// that call to fail; the remaining fields record what the service did.
type EntryRepository struct {
	GetByDayErr        error
	SaveErr            error
	DeleteErr          error
	CurrentWeightKgErr error

	SaveCalls  int
	Saved      *water.Entry
	DeletedID  string
	DeletedFor string

	RequestedStart time.Time
	RequestedEnd   time.Time

	// Entries seeds GetByDay; Weights maps a profile id to the weight
	// CurrentWeightKg reports for it.
	Entries []water.Entry
	Weights map[string]float64
}

var _ water.Repository = (*EntryRepository)(nil)

// NewEntryRepository builds an empty in-memory repository.
func NewEntryRepository() *EntryRepository {
	return &EntryRepository{Weights: map[string]float64{}}
}

// SeedWeight makes CurrentWeightKg report a weight for a profile.
func (r *EntryRepository) SeedWeight(profileID string, weight float64) {
	r.Weights[profileID] = weight
}

// GetByDay returns the seeded entries, or nil when none were seeded — the same
// shape the SQL implementation returns for an empty day.
func (r *EntryRepository) GetByDay(
	_ context.Context,
	_ string,
	start, end time.Time,
) ([]water.Entry, error) {
	if r.GetByDayErr != nil {
		return nil, r.GetByDayErr
	}

	r.RequestedStart = start
	r.RequestedEnd = end

	if len(r.Entries) == 0 {
		return nil, nil
	}

	out := make([]water.Entry, len(r.Entries))
	copy(out, r.Entries)

	return out, nil
}

// Save records the entry.
func (r *EntryRepository) Save(_ context.Context, e *water.Entry) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}

	stored := *e
	r.SaveCalls++
	r.Saved = &stored

	return nil
}

// Delete records the removal.
func (r *EntryRepository) Delete(_ context.Context, id, profileID string) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}

	r.DeletedID = id
	r.DeletedFor = profileID

	return nil
}

// CurrentWeightKg returns the seeded weight of a profile.
func (r *EntryRepository) CurrentWeightKg(_ context.Context, profileID string) (float64, bool, error) {
	if r.CurrentWeightKgErr != nil {
		return 0, false, r.CurrentWeightKgErr
	}

	weight, ok := r.Weights[profileID]
	if !ok {
		return 0, false, nil
	}

	return weight, true, nil
}
