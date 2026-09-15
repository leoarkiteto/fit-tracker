package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/aiplanning"
)

func TestAccept_StoresThePlanForTheProfile(t *testing.T) {
	h := newTestHarness(t)

	planned := []aiplanning.PlannedWorkout{
		{Name: "Push", Goal: "hypertrophy", Days: []string{"monday", "thursday"}},
		{Name: "Pull", Goal: "hypertrophy", Days: []string{"tuesday"}},
	}

	if err := h.service.Accept(context.Background(), testProfileID, planned); err != nil {
		t.Fatalf("Accept() error = %v, want nil", err)
	}

	if h.writer.CreateAllCalls != 1 {
		t.Errorf("CreateAllCalls = %d, want 1", h.writer.CreateAllCalls)
	}
	if h.writer.ProfileID != testProfileID {
		t.Errorf("ProfileID = %q, want %q", h.writer.ProfileID, testProfileID)
	}
	if len(h.writer.Planned) != 2 {
		t.Fatalf("len(Planned) = %d, want 2", len(h.writer.Planned))
	}
	if h.writer.Planned[0].Name != "Push" {
		t.Errorf("Planned[0].Name = %q, want %q", h.writer.Planned[0].Name, "Push")
	}
}

func TestAccept_EmptyPlanStillSucceeds(t *testing.T) {
	h := newTestHarness(t)

	if err := h.service.Accept(context.Background(), testProfileID, nil); err != nil {
		t.Fatalf("Accept() error = %v, want nil", err)
	}
}

func TestAccept_WriterFailure(t *testing.T) {
	h := newTestHarness(t)
	h.writer.CreateAllErr = errors.New("database is gone")

	err := h.service.Accept(context.Background(), testProfileID, nil)

	if appErr := requireAppError(t, err); appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
