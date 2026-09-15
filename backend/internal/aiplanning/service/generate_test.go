package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"fittracker-api/internal/aiplanning"
)

func TestGenerate_WrapsTheModelAnswerWithServerFields(t *testing.T) {
	h := newTestHarness(t)

	age := 31
	weight := 78.5
	h.profiles.Profiles[testProfileID] = aiplanning.Profile{
		Name:          "Leo",
		Age:           &age,
		CurrentWeight: &weight,
	}
	h.planner.Response = `{"summary":"4 training days","rationale":"because","workouts":[{"name":"Push"}]}`

	plan, err := h.service.Generate(context.Background(), testProfileID, "hypertrophy")
	if err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}

	if plan.PlanID == "" {
		t.Error("PlanID = empty, want a generated id")
	}
	if plan.Summary != "4 training days" {
		t.Errorf("Summary = %v, want %q", plan.Summary, "4 training days")
	}
	if plan.Rationale != "because" {
		t.Errorf("Rationale = %v, want %q", plan.Rationale, "because")
	}
	if plan.Workouts == nil {
		t.Error("Workouts = nil, want the model's workouts")
	}
	if _, err := time.Parse(time.RFC3339, plan.GeneratedAt); err != nil {
		t.Errorf("GeneratedAt = %q, want RFC3339, got error %v", plan.GeneratedAt, err)
	}

	// The prompt must carry real values, never pointer addresses.
	for _, want := range []string{"User: Leo", "Age: 31", "Weight: 78.5kg", "Goal: hypertrophy"} {
		if !strings.Contains(h.planner.LastUserPrompt, want) {
			t.Errorf("prompt = %q, want it to contain %q", h.planner.LastUserPrompt, want)
		}
	}
	if h.planner.LastSystemPrompt != systemPrompt {
		t.Errorf("system prompt = %q, want %q", h.planner.LastSystemPrompt, systemPrompt)
	}
}

func TestGenerate_MissingAgeAndWeightReadAsUnknown(t *testing.T) {
	h := newTestHarness(t)
	h.profiles.Profiles[testProfileID] = aiplanning.Profile{Name: "Leo"}

	if _, err := h.service.Generate(context.Background(), testProfileID, "strength"); err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}

	for _, want := range []string{"Age: unknown", "Weight: unknownkg"} {
		if !strings.Contains(h.planner.LastUserPrompt, want) {
			t.Errorf("prompt = %q, want it to contain %q", h.planner.LastUserPrompt, want)
		}
	}
}

func TestGenerate_UnknownProfile(t *testing.T) {
	h := newTestHarness(t)

	_, err := h.service.Generate(context.Background(), "missing", "hypertrophy")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr.Message != "Profile not found" {
		t.Errorf("message = %q, want %q", appErr.Message, "Profile not found")
	}
	if h.planner.PlanCalls != 0 {
		t.Errorf("PlanCalls = %d, want 0 — the model must not be called", h.planner.PlanCalls)
	}
}

func TestGenerate_ProfileLookupFailure(t *testing.T) {
	h := newTestHarness(t)
	h.profiles.GetProfileErr = errors.New("database is gone")

	_, err := h.service.Generate(context.Background(), testProfileID, "hypertrophy")

	if appErr := requireAppError(t, err); appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}

func TestGenerate_PlannerUnavailable(t *testing.T) {
	h := newTestHarness(t)
	h.profiles.Profiles[testProfileID] = aiplanning.Profile{Name: "Leo"}
	h.planner.PlanErr = errors.New("connection refused")

	_, err := h.service.Generate(context.Background(), testProfileID, "hypertrophy")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusServiceUnavailable)
	}
	if !strings.Contains(appErr.Message, "Ollama is not reachable") {
		t.Errorf("message = %q, want the Ollama reachability hint", appErr.Message)
	}
}

func TestGenerate_InvalidJSONFromModel(t *testing.T) {
	h := newTestHarness(t)
	h.profiles.Profiles[testProfileID] = aiplanning.Profile{Name: "Leo"}
	h.planner.Response = "sure! here is your plan:"

	_, err := h.service.Generate(context.Background(), testProfileID, "hypertrophy")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
	if !strings.Contains(appErr.Message, "invalid JSON") {
		t.Errorf("message = %q, want the invalid-JSON hint", appErr.Message)
	}
}
