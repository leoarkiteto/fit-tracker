package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"fittracker-api/internal/workout/mock"
	"fittracker-api/pkg/apperror"
)

func TestGetToday_MatchesLowercaseWeekday(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.SeedToday(testProfileID, sampleWorkout())
	svc := newTestService(t, repo)

	workouts, err := svc.GetToday(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetToday() error = %v, want nil", err)
	}

	want := strings.ToLower(time.Now().Weekday().String())
	if repo.RequestedDay != want {
		t.Errorf("requested day = %q, want %q", repo.RequestedDay, want)
	}
	if len(workouts) != 1 {
		t.Fatalf("len(workouts) = %d, want 1", len(workouts))
	}
	// Unlike the other read endpoints this one has never echoed the profile id.
	if workouts[0].UserProfileID != "" {
		t.Errorf("UserProfileID = %q, want the empty string", workouts[0].UserProfileID)
	}
}

func TestGetToday_RepositoryFailure(t *testing.T) {
	repo := mock.NewWorkoutRepository()
	repo.GetByProfileAndDayErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetToday(context.Background(), testProfileID)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusInternalServerError {
		t.Fatalf("error = %v, want a 500 apperror", err)
	}
}
