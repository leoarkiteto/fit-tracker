package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestUpdate_UsesPathIDAndRegeneratesExerciseIDs(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.Seed(sampleWorkout())
	svc := newTestService(t, repo)

	submitted := sampleWorkout()
	submitted.Name = "Push day v2"
	submitted.ID = ""

	if err := svc.Update(context.Background(), "workout-1", testProfileID, &submitted); err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	if repo.UpdateCalls != 1 {
		t.Errorf("UpdateCalls = %d, want 1", repo.UpdateCalls)
	}
	if repo.UpdatedID != "workout-1" || repo.UpdatedProfileID != testProfileID {
		t.Errorf(
			"repository got (id=%q, profile=%q), want (%q, %q)",
			repo.UpdatedID,
			repo.UpdatedProfileID,
			"workout-1",
			testProfileID,
		)
	}
	if submitted.ID != "workout-1" {
		t.Errorf("ID = %q, want the id from the path", submitted.ID)
	}
	if len(submitted.Exercises) != 1 {
		t.Fatalf("len(Exercises) = %d, want 1", len(submitted.Exercises))
	}
	// Exercises are re-created on every update, so the id is always fresh.
	if submitted.Exercises[0].ID == "" || submitted.Exercises[0].ID == "exercise-1" {
		t.Errorf("exercise id = %q, want a newly generated id", submitted.Exercises[0].ID)
	}
}

func TestUpdate_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.UpdateErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	submitted := sampleWorkout()

	err := svc.Update(context.Background(), "workout-1", testProfileID, &submitted)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
