package service

import (
	"errors"
	"testing"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/pkg/apperror"
)

// testProfileID is the profile every bioimpedance test scopes to.
const testProfileID = "profile-1"

// newTestService builds the bioimpedance service on an in-memory repository.
func newTestService(t *testing.T, repo bioimpedance.Repository) *Service {
	t.Helper()
	return New(repo)
}

// requireAppError asserts that err carries an HTTP status and returns it.
func requireAppError(t *testing.T, err error) *apperror.Error {
	t.Helper()

	if err == nil {
		t.Fatal("error = nil, want an *apperror.Error")
	}

	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v (%T), want *apperror.Error", err, err)
	}

	return appErr
}

// requireStatus asserts the exact HTTP status an error carries.
func requireStatus(t *testing.T, err error, want int) {
	t.Helper()

	if appErr := requireAppError(t, err); appErr.Status != want {
		t.Errorf("status = %d, want %d", appErr.Status, want)
	}
}
