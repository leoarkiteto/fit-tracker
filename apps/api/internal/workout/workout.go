// Package workout owns the workout and exercise domain.
package workout

import (
	"time"

	"fittracker-api/pkg/apperror"
)

// Exercise is a single movement inside a workout.
type Exercise struct {
	Weight      *float64 `json:"weight"`
	Notes       *string  `json:"notes"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MuscleGroup string   `json:"muscleGroup"`
	WorkoutID   string   `json:"workoutId"`
	Sets        int      `json:"sets"`
	Reps        int      `json:"reps"`
	RestSeconds int      `json:"restSeconds"`
}

// Workout is a planned training session.
type Workout struct {
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Description   *string    `json:"description"`
	CompletedAt   *time.Time `json:"completedAt"`
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Goal          string     `json:"goal"`
	UserProfileID string     `json:"userProfileId"`
	Days          []string   `json:"days"`
	Exercises     []Exercise `json:"exercises"`
	IsCompleted   bool       `json:"isCompleted"`
}

// ErrWorkoutNotFound is returned when a workout id does not exist for the
// profile in the URL.
var ErrWorkoutNotFound = apperror.NotFound("Workout not found")
