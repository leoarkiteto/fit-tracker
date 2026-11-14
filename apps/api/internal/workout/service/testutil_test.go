package service

import (
	"testing"

	"fittracker-api/internal/workout"
)

// testProfileID is the profile every workout test scopes to.
const testProfileID = "profile-1"

// newTestService builds the workout service on an in-memory repository.
func newTestService(t *testing.T, repo workout.Repository) *Service {
	t.Helper()
	return New(repo)
}

// sampleWorkout returns a workout carrying one exercise.
func sampleWorkout() workout.Workout {
	weight := 40.0
	notes := "keep the elbows in"

	return workout.Workout{
		ID:   "workout-1",
		Name: "Push day",
		Goal: "hypertrophy",
		Days: []string{"monday", "thursday"},
		Exercises: []workout.Exercise{
			{
				ID:          "exercise-1",
				Name:        "Bench press",
				MuscleGroup: "chest",
				Sets:        4,
				Reps:        8,
				Weight:      &weight,
				RestSeconds: 90,
				Notes:       &notes,
			},
		},
	}
}
