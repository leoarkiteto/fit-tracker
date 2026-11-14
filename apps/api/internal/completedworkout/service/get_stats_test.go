package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/completedworkout/mock"
)

func TestGetStats_ComputesMinutesFromSeconds(t *testing.T) {
	repo := mock.NewCompletedWorkoutRepository()
	repo.Totals[testProfileID] = 12
	repo.ThisWeek[testProfileID] = 3
	repo.Seconds[testProfileID] = 3671
	svc := newTestService(t, repo)

	before := time.Now()
	stats, err := svc.GetStats(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetStats() error = %v, want nil", err)
	}

	if stats.TotalWorkoutsCompleted != 12 {
		t.Errorf("TotalWorkoutsCompleted = %d, want 12", stats.TotalWorkoutsCompleted)
	}
	if stats.WorkoutsThisWeek != 3 {
		t.Errorf("WorkoutsThisWeek = %d, want 3", stats.WorkoutsThisWeek)
	}
	// Integer division, exactly as before: 3671 / 60 == 61.
	if stats.TotalMinutesSpent != 61 {
		t.Errorf("TotalMinutesSpent = %d, want 61", stats.TotalMinutesSpent)
	}

	// "This week" is the last seven days.
	wanted := before.Add(-statsWindow)
	if repo.RequestedSince.Before(wanted.Add(-time.Second)) ||
		repo.RequestedSince.After(wanted.Add(time.Second)) {
		t.Errorf("RequestedSince = %v, want about %v", repo.RequestedSince, wanted)
	}
}

func TestGetStats_ZeroForEmptyHistory(t *testing.T) {
	svc := newTestService(t, mock.NewCompletedWorkoutRepository())

	stats, err := svc.GetStats(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetStats() error = %v, want nil", err)
	}

	if stats.TotalWorkoutsCompleted != 0 || stats.WorkoutsThisWeek != 0 || stats.TotalMinutesSpent != 0 {
		t.Errorf("stats = %+v, want every counter at zero", stats)
	}
}

func TestGetStats_RepositoryFailures(t *testing.T) {
	boom := errors.New("database is gone")

	cases := map[string]func(*mock.CompletedWorkoutRepository){
		"count":     func(r *mock.CompletedWorkoutRepository) { r.CountByProfileErr = boom },
		"this week": func(r *mock.CompletedWorkoutRepository) { r.CountByProfileSinceErr = boom },
		"summing":   func(r *mock.CompletedWorkoutRepository) { r.SumDurationByProfileErr = boom },
	}

	for name, breakRepo := range cases {
		t.Run(name, func(t *testing.T) {
			repo := mock.NewCompletedWorkoutRepository()
			breakRepo(repo)
			svc := newTestService(t, repo)

			_, err := svc.GetStats(context.Background(), testProfileID)

			requireStatus(t, err, http.StatusInternalServerError)
		})
	}
}
