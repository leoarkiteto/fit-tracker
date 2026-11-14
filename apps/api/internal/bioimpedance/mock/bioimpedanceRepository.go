// Package mock holds hand-written test doubles for the bioimpedance slice.
package mock

import (
	"context"

	"fittracker-api/internal/bioimpedance"
)

// MeasurementRepository is an in-memory bioimpedance.Repository. Set any Err
// field to force that call to fail; the remaining fields record what the
// service did.
type MeasurementRepository struct {
	GetByProfileErr error
	GetLatestErr    error
	SaveErr         error
	DeleteErr       error

	SaveCalls  int
	Saved      *bioimpedance.Measurement
	DeletedID  string
	DeletedFor string

	// History seeds GetByProfile; Latest seeds GetLatest.
	History []bioimpedance.Measurement
	Latest  *bioimpedance.Measurement
}

var _ bioimpedance.Repository = (*MeasurementRepository)(nil)

// NewMeasurementRepository builds an empty in-memory repository.
func NewMeasurementRepository() *MeasurementRepository {
	return &MeasurementRepository{}
}

// GetByProfile returns the seeded history, or nil when none was seeded — the
// same shape the SQL implementation returns for an empty table.
func (r *MeasurementRepository) GetByProfile(
	_ context.Context,
	_ string,
) ([]bioimpedance.Measurement, error) {
	if r.GetByProfileErr != nil {
		return nil, r.GetByProfileErr
	}

	if len(r.History) == 0 {
		return nil, nil
	}

	out := make([]bioimpedance.Measurement, len(r.History))
	copy(out, r.History)

	return out, nil
}

// GetLatest returns the seeded latest measurement.
func (r *MeasurementRepository) GetLatest(
	_ context.Context,
	_ string,
) (*bioimpedance.Measurement, bool, error) {
	if r.GetLatestErr != nil {
		return nil, false, r.GetLatestErr
	}

	if r.Latest == nil {
		return nil, false, nil
	}

	stored := *r.Latest
	return &stored, true, nil
}

// Save records the measurement.
func (r *MeasurementRepository) Save(_ context.Context, m *bioimpedance.Measurement) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}

	stored := *m
	r.SaveCalls++
	r.Saved = &stored

	return nil
}

// Delete records the removal.
func (r *MeasurementRepository) Delete(_ context.Context, id, profileID string) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}

	r.DeletedID = id
	r.DeletedFor = profileID

	return nil
}
