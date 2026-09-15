package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/completedworkout/mock"
)

func TestComplete_StoresEntryForProfile(t *testing.T) {
	repo := mock.NewCompletedWorkoutRepository()
	svc := newTestService(t, repo)

	before := time.Now()
	entry, err := svc.Complete(context.Background(), testProfileID, "workout-1", 1845)
	if err != nil {
		t.Fatalf("Complete() error = %v, want nil", err)
	}

	if repo.SaveCalls != 1 {
		t.Errorf("SaveCalls = %d, want 1", repo.SaveCalls)
	}
	if entry.ID == "" {
		t.Error("ID = empty, want a generated id")
	}
	if entry.WorkoutID != "workout-1" {
		t.Errorf("WorkoutID = %q, want %q", entry.WorkoutID, "workout-1")
	}
	if entry.UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", entry.UserProfileID, testProfileID)
	}
	if entry.DurationSeconds != 1845 {
		t.Errorf("DurationSeconds = %d, want 1845", entry.DurationSeconds)
	}
	if entry.CompletedAt.Before(before) || entry.CompletedAt.After(time.Now()) {
		t.Errorf("CompletedAt = %v, want the moment of the call", entry.CompletedAt)
	}
	if repo.Saved == nil || repo.Saved.ID != entry.ID {
		t.Error("the repository did not receive the created entry")
	}
}

func TestComplete_RepositoryFailure(t *testing.T) {
	repo := mock.NewCompletedWorkoutRepository()
	repo.SaveErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.Complete(context.Background(), testProfileID, "workout-1", 1845)

	requireStatus(t, err, http.StatusInternalServerError)
}
