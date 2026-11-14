package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestDelete_ScopesToProfile(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.Seed(sampleWorkout())
	svc := newTestService(t, repo)

	if err := svc.Delete(context.Background(), "workout-1", testProfileID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	if repo.DeleteCalls != 1 {
		t.Errorf("DeleteCalls = %d, want 1", repo.DeleteCalls)
	}
	if repo.DeletedID != "workout-1" || repo.DeletedProfileID != testProfileID {
		t.Errorf(
			"repository got (id=%q, profile=%q), want (%q, %q)",
			repo.DeletedID,
			repo.DeletedProfileID,
			"workout-1",
			testProfileID,
		)
	}
}

func TestDelete_UnknownWorkoutIsNotAnError(t *testing.T) {
	svc := newTestService(t, mock.NewWorkoutRepository())

	// The previous implementation answered 204 for an unknown id and this keeps
	// that behaviour.
	if err := svc.Delete(context.Background(), "missing", testProfileID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
}

func TestDelete_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.DeleteErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	err := svc.Delete(context.Background(), "workout-1", testProfileID)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
