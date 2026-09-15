// Package mock holds hand-written test doubles for the workout slice, filling
// the role mockery-generated mocks play in the reference architecture.
package mock

import (
	"context"

	"fittracker-api/internal/workout"
)

// WorkoutRepository is an in-memory workout.Repository. Set any Err field to
// force that call to fail; the fields below record what the service did.
type WorkoutRepository struct {
	GetByProfileErr       error
	GetByProfileAndDayErr error
	GetByIDErr            error
	SaveErr               error
	UpdateErr             error
	DeleteErr             error

	SaveCalls   int
	UpdateCalls int
	DeleteCalls int

	SavedProfileID   string
	SavedWorkout     *workout.Workout
	UpdatedID        string
	UpdatedProfileID string
	UpdatedWorkout   *workout.Workout
	DeletedID        string
	DeletedProfileID string
	RequestedDay     string

	// ByProfile and Today seed the two list queries; ByID seeds GetByID.
	ByProfile map[string][]workout.Workout
	Today     map[string][]workout.Workout
	ByID      map[string]workout.Workout
}

var _ workout.Repository = (*WorkoutRepository)(nil)

// NewWorkoutRepository builds an empty in-memory repository.
func NewWorkoutRepository() *WorkoutRepository {
	return &WorkoutRepository{
		ByProfile: map[string][]workout.Workout{},
		Today:     map[string][]workout.Workout{},
		ByID:      map[string]workout.Workout{},
	}
}

// SeedProfile registers the workouts returned by GetByProfile.
func (r *WorkoutRepository) SeedProfile(profileID string, workouts ...workout.Workout) {
	r.ByProfile[profileID] = clones(workouts)
}

// SeedToday registers the workouts returned by GetByProfileAndDay.
func (r *WorkoutRepository) SeedToday(profileID string, workouts ...workout.Workout) {
	r.Today[profileID] = clones(workouts)
}

// Seed registers the workout returned by GetByID.
func (r *WorkoutRepository) Seed(w workout.Workout) {
	r.ByID[w.ID] = clone(w)
}

// GetByProfile returns the seeded workouts of a profile, or nil when none were
// seeded — the same shape the SQL implementation returns for an empty table.
func (r *WorkoutRepository) GetByProfile(_ context.Context, profileID string) ([]workout.Workout, error) {
	if r.GetByProfileErr != nil {
		return nil, r.GetByProfileErr
	}

	stored, ok := r.ByProfile[profileID]
	if !ok {
		return nil, nil
	}

	return clones(stored), nil
}

// GetByProfileAndDay returns the seeded workouts for the requested weekday.
func (r *WorkoutRepository) GetByProfileAndDay(
	_ context.Context,
	profileID, day string,
) ([]workout.Workout, error) {
	if r.GetByProfileAndDayErr != nil {
		return nil, r.GetByProfileAndDayErr
	}

	r.RequestedDay = day

	stored, ok := r.Today[profileID]
	if !ok {
		return nil, nil
	}

	return clones(stored), nil
}

// GetByID returns the seeded workout.
func (r *WorkoutRepository) GetByID(_ context.Context, id, _ string) (*workout.Workout, bool, error) {
	if r.GetByIDErr != nil {
		return nil, false, r.GetByIDErr
	}

	stored, ok := r.ByID[id]
	if !ok {
		return nil, false, nil
	}

	found := clone(stored)
	return &found, true, nil
}

// Save records the created workout.
func (r *WorkoutRepository) Save(_ context.Context, profileID string, w *workout.Workout) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}

	stored := clone(*w)

	r.SaveCalls++
	r.SavedProfileID = profileID
	r.SavedWorkout = &stored
	r.ByID[w.ID] = stored

	return nil
}

// Update records the replaced workout.
func (r *WorkoutRepository) Update(_ context.Context, id, profileID string, w *workout.Workout) error {
	if r.UpdateErr != nil {
		return r.UpdateErr
	}

	stored := clone(*w)

	r.UpdateCalls++
	r.UpdatedID = id
	r.UpdatedProfileID = profileID
	r.UpdatedWorkout = &stored
	r.ByID[id] = stored

	return nil
}

// Delete records the removed workout.
func (r *WorkoutRepository) Delete(_ context.Context, id, profileID string) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}

	r.DeleteCalls++
	r.DeletedID = id
	r.DeletedProfileID = profileID
	delete(r.ByID, id)

	return nil
}

func clones(workouts []workout.Workout) []workout.Workout {
	out := make([]workout.Workout, 0, len(workouts))
	for _, w := range workouts {
		out = append(out, clone(w))
	}

	return out
}

func clone(w workout.Workout) workout.Workout {
	c := w
	c.Days = append([]string(nil), w.Days...)
	c.Exercises = append([]workout.Exercise(nil), w.Exercises...)

	return c
}
