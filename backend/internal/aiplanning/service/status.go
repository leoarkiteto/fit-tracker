package service

import (
	"context"

	"fittracker-api/internal/aiplanning"
)

// BackendStatus reports whether the model backend is reachable.
func (s *Service) BackendStatus(ctx context.Context) *aiplanning.BackendStatus {
	ctx, cancel := context.WithTimeout(ctx, aiplanning.HealthTimeout)
	defer cancel()

	available := s.planner.Available(ctx) == nil

	provider := "Ollama (via LangChainGo)"
	if !available {
		provider = "unreachable"
	}

	return &aiplanning.BackendStatus{
		Available: available,
		Provider:  provider,
		Model:     s.planner.Model(),
		Endpoint:  s.planner.Endpoint(),
	}
}
