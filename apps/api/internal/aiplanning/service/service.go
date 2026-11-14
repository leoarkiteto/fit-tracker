// Package service holds the AI-planning use cases. It depends only on the
// domain and its ports, so it can be unit tested without HTTP, a database or a
// running Ollama.
package service

import (
	"fittracker-api/internal/aiplanning"
)

// Service holds the AI-planning use cases.
type Service struct {
	profiles aiplanning.ProfileRepository
	writer   aiplanning.WorkoutWriter
	planner  aiplanning.Planner
}

// New builds the AI-planning service.
func New(
	profiles aiplanning.ProfileRepository,
	writer aiplanning.WorkoutWriter,
	planner aiplanning.Planner,
) *Service {
	return &Service{profiles: profiles, writer: writer, planner: planner}
}
