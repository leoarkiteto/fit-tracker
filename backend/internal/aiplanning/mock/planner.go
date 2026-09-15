// Package mock holds hand-written test doubles for the aiplanning slice.
package mock

import (
	"context"

	"fittracker-api/internal/aiplanning"
)

// Planner is a controllable aiplanning.Planner. It records the prompts it was
// given so tests can assert what the model would have seen.
type Planner struct {
	PlanErr      error
	AvailableErr error

	// Response is the raw answer Plan returns.
	Response string

	ModelName   string
	EndpointURL string

	PlanCalls        int
	LastSystemPrompt string
	LastUserPrompt   string
}

var _ aiplanning.Planner = (*Planner)(nil)

// NewPlanner builds a Planner answering with a minimal valid plan.
func NewPlanner() *Planner {
	return &Planner{
		Response:    `{"summary":"a summary","rationale":"a rationale","workouts":[]}`,
		ModelName:   "test-model",
		EndpointURL: "http://test-endpoint",
	}
}

// Plan records the prompts and returns the configured answer.
func (p *Planner) Plan(_ context.Context, systemPrompt, userPrompt string) (string, error) {
	p.PlanCalls++
	p.LastSystemPrompt = systemPrompt
	p.LastUserPrompt = userPrompt

	if p.PlanErr != nil {
		return "", p.PlanErr
	}

	return p.Response, nil
}

// Available reports the configured health error.
func (p *Planner) Available(_ context.Context) error { return p.AvailableErr }

// Model is the configured model name.
func (p *Planner) Model() string { return p.ModelName }

// Endpoint is the configured backend URL.
func (p *Planner) Endpoint() string { return p.EndpointURL }
