package aiplanning

import "context"

// ProfileRepository is the outbound port for reading the profile a plan is
// generated for.
type ProfileRepository interface {
	// GetProfile returns the planning-relevant fields of a profile, with
	// found=false when the id is unknown.
	GetProfile(ctx context.Context, profileID string) (*Profile, bool, error)
}

// WorkoutWriter is the outbound port for storing the workouts of an accepted
// plan.
type WorkoutWriter interface {
	// CreateAll stores every workout of the plan for a profile in a single
	// transaction, generating the ids the plan left out.
	CreateAll(ctx context.Context, profileID string, planned []PlannedWorkout) error
}

// Planner is the outbound port for the language model behind the feature.
type Planner interface {
	// Plan asks the model for a workout plan and returns its raw answer.
	Plan(ctx context.Context, systemPrompt, userPrompt string) (string, error)

	// Available reports whether the model backend answers at all.
	Available(ctx context.Context) error

	// Model is the configured model name.
	Model() string

	// Endpoint is the configured backend URL.
	Endpoint() string
}
