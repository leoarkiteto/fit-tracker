package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestGetAll_EchoesProfileID(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.SeedProfile(testProfileID, sampleWorkout())
	svc := newTestService(t, repo)

	workouts, err := svc.GetAll(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	if len(workouts) != 1 {
		t.Fatalf("len(workouts) = %d, want 1", len(workouts))
	}
	if workouts[0].UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", workouts[0].UserProfileID, testProfileID)
	}
	if len(workouts[0].Exercises) != 1 {
		t.Fatalf("len(Exercises) = %d, want 1", len(workouts[0].Exercises))
	}
	if workouts[0].Exercises[0].Name != "Bench press" {
		t.Errorf("Exercises[0].Name = %q, want %q", workouts[0].Exercises[0].Name, "Bench press")
	}
}

func TestGetAll_EmptyStaysNull(t *testing.T) {
	svc := newTestService(t, mock.NewWorkoutRepository())

	workouts, err := svc.GetAll(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	// A nil slice is what the API has always serialised as JSON null.
	if workouts != nil {
		t.Errorf("workouts = %#v, want nil", workouts)
	}
}

func TestGetAll_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.GetByProfileErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetAll(context.Background(), testProfileID)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
