// Package mock holds hand-written test doubles for the completed-workout slice.
package mock

import (
	"context"
	"time"

	"fittracker-api/internal/completedworkout"
)

// CompletedWorkoutRepository is an in-memory completedworkout.Repository. Set
// any Err field to force that call to fail; the maps seed the queries and the
// remaining fields record what the service did.
type CompletedWorkoutRepository struct {
	GetByProfileErr         error
	SaveErr                 error
	CountByProfileErr       error
	CountByProfileSinceErr  error
	SumDurationByProfileErr error

	SaveCalls      int
	Saved          *completedworkout.CompletedWorkout
	RequestedSince time.Time

	// History seeds GetByProfile; Totals, ThisWeek and Seconds seed the three
	// statistics queries.
	History  map[string][]completedworkout.CompletedWorkout
	Totals   map[string]int
	ThisWeek map[string]int
	Seconds  map[string]int
}

var _ completedworkout.Repository = (*CompletedWorkoutRepository)(nil)

// NewCompletedWorkoutRepository builds an empty in-memory repository.
func NewCompletedWorkoutRepository() *CompletedWorkoutRepository {
	return &CompletedWorkoutRepository{
		History:  map[string][]completedworkout.CompletedWorkout{},
		Totals:   map[string]int{},
		ThisWeek: map[string]int{},
		Seconds:  map[string]int{},
	}
}

// GetByProfile returns the seeded history, or nil when none was seeded — the
// same shape the SQL implementation returns for an empty table.
func (r *CompletedWorkoutRepository) GetByProfile(
	_ context.Context,
	profileID string,
) ([]completedworkout.CompletedWorkout, error) {
	if r.GetByProfileErr != nil {
		return nil, r.GetByProfileErr
	}

	stored, ok := r.History[profileID]
	if !ok {
		return nil, nil
	}

	out := make([]completedworkout.CompletedWorkout, len(stored))
	copy(out, stored)

	return out, nil
}

// Save records the completed workout.
func (r *CompletedWorkoutRepository) Save(_ context.Context, c *completedworkout.CompletedWorkout) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}

	stored := *c
	r.SaveCalls++
	r.Saved = &stored

	return nil
}

// CountByProfile returns the seeded total.
func (r *CompletedWorkoutRepository) CountByProfile(_ context.Context, profileID string) (int, error) {
	if r.CountByProfileErr != nil {
		return 0, r.CountByProfileErr
	}

	return r.Totals[profileID], nil
}

// CountByProfileSince returns the seeded weekly count.
func (r *CompletedWorkoutRepository) CountByProfileSince(
	_ context.Context,
	profileID string,
	since time.Time,
) (int, error) {
	if r.CountByProfileSinceErr != nil {
		return 0, r.CountByProfileSinceErr
	}

	r.RequestedSince = since

	return r.ThisWeek[profileID], nil
}

// SumDurationByProfile returns the seeded number of seconds.
func (r *CompletedWorkoutRepository) SumDurationByProfile(_ context.Context, profileID string) (int, error) {
	if r.SumDurationByProfileErr != nil {
		return 0, r.SumDurationByProfileErr
	}

	return r.Seconds[profileID], nil
}
