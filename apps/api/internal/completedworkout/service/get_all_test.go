package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/completedworkout"
	"fittracker-api/internal/completedworkout/mock"
)

func TestGetAll_EchoesProfileID(t *testing.T) {
	repo := mock.NewCompletedWorkoutRepository()
	repo.History[testProfileID] = []completedworkout.CompletedWorkout{
		{ID: "completed-1", WorkoutID: "workout-1", DurationSeconds: 1800},
	}
	svc := newTestService(t, repo)

	history, err := svc.GetAll(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	if len(history) != 1 {
		t.Fatalf("len(history) = %d, want 1", len(history))
	}
	if history[0].UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", history[0].UserProfileID, testProfileID)
	}
	if history[0].WorkoutID != "workout-1" {
		t.Errorf("WorkoutID = %q, want %q", history[0].WorkoutID, "workout-1")
	}
}

func TestGetAll_EmptyStaysNull(t *testing.T) {
	svc := newTestService(t, mock.NewCompletedWorkoutRepository())

	history, err := svc.GetAll(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	// A nil slice is what the API has always serialised as JSON null.
	if history != nil {
		t.Errorf("history = %#v, want nil", history)
	}
}

func TestGetAll_RepositoryFailure(t *testing.T) {
	repo := mock.NewCompletedWorkoutRepository()
	repo.GetByProfileErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetAll(context.Background(), testProfileID)

	requireStatus(t, err, http.StatusInternalServerError)
}
