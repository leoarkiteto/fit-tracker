package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestGetByID_Success(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.Seed(sampleWorkout())
	svc := newTestService(t, repo)

	result, err := svc.GetByID(context.Background(), "workout-1", testProfileID)
	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}

	if result.Name != "Push day" {
		t.Errorf("Name = %q, want %q", result.Name, "Push day")
	}
	if result.UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", result.UserProfileID, testProfileID)
	}
	if len(result.Days) != 2 || result.Days[0] != "monday" {
		t.Errorf("Days = %v, want [monday thursday]", result.Days)
	}
	if len(result.Exercises) != 1 {
		t.Fatalf("len(Exercises) = %d, want 1", len(result.Exercises))
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := newTestService(t, mock.NewWorkoutRepository())

	_, err := svc.GetByID(context.Background(), "missing", testProfileID)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *apperror.Error", err)
	}
	if appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr.Message != "Workout not found" {
		t.Errorf("message = %q, want %q", appErr.Message, "Workout not found")
	}
}

func TestGetByID_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.GetByIDErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetByID(context.Background(), "workout-1", testProfileID)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
