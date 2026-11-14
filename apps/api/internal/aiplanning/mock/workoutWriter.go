package mock

import (
	"context"

	"fittracker-api/internal/aiplanning"
)

// WorkoutWriter records the plans handed to it.
type WorkoutWriter struct {
	CreateAllErr error

	CreateAllCalls int
	ProfileID      string
	Planned        []aiplanning.PlannedWorkout
}

var _ aiplanning.WorkoutWriter = (*WorkoutWriter)(nil)

// NewWorkoutWriter builds an empty recorder.
func NewWorkoutWriter() *WorkoutWriter { return &WorkoutWriter{} }

// CreateAll records the call.
func (w *WorkoutWriter) CreateAll(
	_ context.Context,
	profileID string,
	planned []aiplanning.PlannedWorkout,
) error {
	if w.CreateAllErr != nil {
		return w.CreateAllErr
	}

	w.CreateAllCalls++
	w.ProfileID = profileID
	w.Planned = append([]aiplanning.PlannedWorkout(nil), planned...)

	return nil
}
