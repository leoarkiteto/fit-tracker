package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/workout"
	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestCreate_GeneratesIDsAndScopesToProfile(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	svc := newTestService(t, repo)

	created := workout.Workout{
		Name:      "Leg day",
		Goal:      "strength",
		Days:      []string{"tuesday"},
		Exercises: []workout.Exercise{{Name: "Squat", MuscleGroup: "legs", Sets: 5, Reps: 5}},
	}

	if err := svc.Create(context.Background(), testProfileID, &created); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if repo.SaveCalls != 1 {
		t.Errorf("SaveCalls = %d, want 1", repo.SaveCalls)
	}
	if repo.SavedProfileID != testProfileID {
		t.Errorf("SavedProfileID = %q, want %q", repo.SavedProfileID, testProfileID)
	}
	if created.ID == "" {
		t.Error("ID = empty, want a generated workout id")
	}
	if len(created.Exercises) != 1 || created.Exercises[0].ID == "" {
		t.Error("exercise id was not generated")
	}
	if repo.SavedWorkout == nil || repo.SavedWorkout.ID != created.ID {
		t.Error("the repository did not receive the generated workout")
	}

	// The response body echoes the submitted body, so the profile id from the
	// path must not leak into it.
	if created.UserProfileID != "" {
		t.Errorf("UserProfileID = %q, want the empty string", created.UserProfileID)
	}
}

func TestCreate_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.SaveErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	created := sampleWorkout()

	err := svc.Create(context.Background(), testProfileID, &created)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
