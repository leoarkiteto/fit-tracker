// Package aiplanning owns the AI workout-planning domain: it asks a local
// Ollama model for a weekly plan and stores the plans the user accepts.
package aiplanning

import (
	"net/http"
	"time"

	"fittracker-api/pkg/apperror"
)

// Profile is the slice of a trainee profile a plan is generated from.
type Profile struct {
	Age           *int
	CurrentWeight *float64
	Name          string
}

// GeneratedPlan is the answer of POST /api/ai/planning/generate. Summary,
// Rationale and Workouts are passed through from the model untouched, so they
// stay untyped.
type GeneratedPlan struct {
	PlanID      string `json:"planId"`
	Summary     any    `json:"summary"`
	Rationale   any    `json:"rationale"`
	Workouts    any    `json:"workouts"`
	GeneratedAt string `json:"generatedAt"`
}

// BackendStatus describes the planning backend, as reported by
// GET /api/ai/planning/status.
type BackendStatus struct {
	Available bool   `json:"available"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Endpoint  string `json:"endpoint"`
}

// PlannedWorkout is a workout as it arrives in an accepted plan. Only the
// fields the API persists are decoded.
type PlannedWorkout struct {
	Description *string           `json:"description"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Goal        string            `json:"goal"`
	Days        []string          `json:"days"`
	Exercises   []PlannedExercise `json:"exercises"`
}

// PlannedExercise is an exercise as it arrives in an accepted plan.
type PlannedExercise struct {
	Weight      *float64 `json:"weight"`
	Notes       *string  `json:"notes"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MuscleGroup string   `json:"muscleGroup"`
	Sets        int      `json:"sets"`
	Reps        int      `json:"reps"`
	RestSeconds int      `json:"restSeconds"`
}

// GenerationTimeout bounds how long a single plan generation may take.
const GenerationTimeout = 60 * time.Second

// HealthTimeout bounds the backend status probe.
const HealthTimeout = 5 * time.Second

// Domain errors. The messages are part of the API contract.
var (
	// ErrProfileNotFound is returned when the profile to plan for is unknown.
	ErrProfileNotFound = apperror.NotFound("Profile not found")
	// ErrPlannerUnavailable is returned when the model backend cannot be reached.
	ErrPlannerUnavailable = apperror.ServiceUnavailable(
		"Ollama is not reachable or model not loaded. Make sure Ollama is running (ollama serve) and the model is pulled (ollama pull llama3.2).",
	)
	// ErrInvalidPlan is returned when the model does not answer with JSON.
	ErrInvalidPlan = apperror.New(
		http.StatusInternalServerError,
		"AI returned invalid JSON. The model may still be loading. Please try again.",
	)
)
