package service

import (
	"errors"
	"testing"

	"fittracker-api/internal/aiplanning/mock"
	"fittracker-api/pkg/apperror"
)

// testProfileID is the profile every AI-planning test scopes to.
const testProfileID = "profile-1"

// testHarness wires the service under test with its three fakes.
type testHarness struct {
	service  *Service
	profiles *mock.ProfileRepository
	writer   *mock.WorkoutWriter
	planner  *mock.Planner
}

// newTestHarness builds a service backed entirely by in-memory doubles, so no
// database and no Ollama are involved.
func newTestHarness(t *testing.T) *testHarness {
	t.Helper()

	profiles := mock.NewProfileRepository()
	writer := mock.NewWorkoutWriter()
	planner := mock.NewPlanner()

	return &testHarness{
		service:  New(profiles, writer, planner),
		profiles: profiles,
		writer:   writer,
		planner:  planner,
	}
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
