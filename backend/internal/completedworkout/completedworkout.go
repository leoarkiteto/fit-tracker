// Package completedworkout owns the training-history domain.
package completedworkout

import "time"

// CompletedWorkout is one finished training session.
type CompletedWorkout struct {
	CompletedAt     time.Time `json:"completedAt"`
	ID              string    `json:"id"`
	WorkoutID       string    `json:"workoutId"`
	UserProfileID   string    `json:"userProfileId"`
	DurationSeconds int       `json:"durationSeconds"`
}

// Stats summarises the history of a profile.
type Stats struct {
	TotalWorkoutsCompleted int `json:"totalWorkoutsCompleted"`
	WorkoutsThisWeek       int `json:"workoutsThisWeek"`
	TotalMinutesSpent      int `json:"totalMinutesSpent"`
}
