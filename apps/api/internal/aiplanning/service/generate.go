package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/aiplanning"
	"fittracker-api/pkg/apperror"
)

const systemPrompt = `You are an experienced fitness coach. Respond ONLY with valid JSON.`

const promptTemplate = `Based on user profile, create a personalized weekly workout plan.
User: %s, Age: %v, Weight: %vkg, Goal: %s.
Generate JSON structure with: summary, rationale, workouts (name, description, goal, days, exercises).`

// unknownField replaces a missing age or weight in the prompt.
const unknownField = "unknown"

// Generate asks the model for a weekly plan for a profile.
func (s *Service) Generate(
	ctx context.Context,
	userProfileID, goal string,
) (*aiplanning.GeneratedPlan, error) {
	profile, found, err := s.profiles.GetProfile(ctx, userProfileID)
	if err != nil {
		log.Printf("ai planning: loading profile %s: %v", userProfileID, err)
		return nil, apperror.Internal()
	}
	if !found {
		return nil, aiplanning.ErrProfileNotFound
	}

	// The prompt needs values, not pointers: a nil age or weight reads as
	// "unknown" so the model is not handed an address.
	age := any(unknownField)
	if profile.Age != nil {
		age = *profile.Age
	}
	weight := any(unknownField)
	if profile.CurrentWeight != nil {
		weight = *profile.CurrentWeight
	}

	userPrompt := fmt.Sprintf(promptTemplate, profile.Name, age, weight, goal)

	// Bounded so a stalled model cannot hold the request open forever.
	ctx, cancel := context.WithTimeout(ctx, aiplanning.GenerationTimeout)
	defer cancel()

	responseText, err := s.planner.Plan(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("ai planning: generating plan: %v", err)
		return nil, aiplanning.ErrPlannerUnavailable
	}

	var modelOutput map[string]any
	if err := json.Unmarshal([]byte(responseText), &modelOutput); err != nil {
		log.Printf("ai planning: model returned invalid JSON: %s, error: %v", responseText, err)
		return nil, aiplanning.ErrInvalidPlan
	}

	return &aiplanning.GeneratedPlan{
		PlanID:      uuid.New().String(),
		Summary:     modelOutput["summary"],
		Rationale:   modelOutput["rationale"],
		Workouts:    modelOutput["workouts"],
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
