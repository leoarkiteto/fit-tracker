package service

import (
	"context"
	"errors"
	"testing"
)

func TestBackendStatus_Available(t *testing.T) {
	h := newTestHarness(t)

	status := h.service.BackendStatus(context.Background())

	if !status.Available {
		t.Error("Available = false, want true")
	}
	if status.Provider != "Ollama (via LangChainGo)" {
		t.Errorf("Provider = %q, want %q", status.Provider, "Ollama (via LangChainGo)")
	}
	if status.Model != "test-model" {
		t.Errorf("Model = %q, want %q", status.Model, "test-model")
	}
	if status.Endpoint != "http://test-endpoint" {
		t.Errorf("Endpoint = %q, want %q", status.Endpoint, "http://test-endpoint")
	}
}

func TestBackendStatus_Unreachable(t *testing.T) {
	h := newTestHarness(t)
	h.planner.AvailableErr = errors.New("connection refused")

	status := h.service.BackendStatus(context.Background())

	if status.Available {
		t.Error("Available = true, want false")
	}
	if status.Provider != "unreachable" {
		t.Errorf("Provider = %q, want %q", status.Provider, "unreachable")
	}
	// The model and endpoint stay reported even when the backend is down.
	if status.Model != "test-model" || status.Endpoint != "http://test-endpoint" {
		t.Errorf("Model/Endpoint = %q/%q, want them still reported", status.Model, status.Endpoint)
	}
}
