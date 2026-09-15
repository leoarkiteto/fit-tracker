package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/water"
	"fittracker-api/internal/water/mock"
)

func TestGetDaily_SumsEntriesAndDerivesGoal(t *testing.T) {
	repo := mock.NewEntryRepository()
	repo.Entries = []water.Entry{
		{ID: "entry-1", AmountMl: 250},
		{ID: "entry-2", AmountMl: 500},
	}
	repo.SeedWeight(testProfileID, 80)
	svc := newTestService(t, repo)

	summary, err := svc.GetDaily(context.Background(), testProfileID, "2026-04-01")
	if err != nil {
		t.Fatalf("GetDaily() error = %v, want nil", err)
	}

	if summary.Date != "2026-04-01" {
		t.Errorf("Date = %q, want %q", summary.Date, "2026-04-01")
	}
	if summary.TotalMl != 750 {
		t.Errorf("TotalMl = %d, want 750", summary.TotalMl)
	}
	// 35 ml per kg of the recorded 80 kg.
	if summary.GoalMl != 2800 {
		t.Errorf("GoalMl = %d, want 2800", summary.GoalMl)
	}
	if len(summary.Entries) != 2 {
		t.Errorf("len(Entries) = %d, want 2", len(summary.Entries))
	}

	// The day window runs from UTC midnight to the next UTC midnight.
	wantStart := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	if !repo.RequestedStart.Equal(wantStart) {
		t.Errorf("RequestedStart = %v, want %v", repo.RequestedStart, wantStart)
	}
	if !repo.RequestedEnd.Equal(wantStart.AddDate(0, 0, 1)) {
		t.Errorf("RequestedEnd = %v, want %v", repo.RequestedEnd, wantStart.AddDate(0, 0, 1))
	}
}

func TestGetDaily_DefaultsToTheDefaultWeightGoalWithoutWeight(t *testing.T) {
	svc := newTestService(t, mock.NewEntryRepository())

	summary, err := svc.GetDaily(context.Background(), testProfileID, "2026-04-01")
	if err != nil {
		t.Fatalf("GetDaily() error = %v, want nil", err)
	}

	// The default 70 kg fallback still applies to a profile with no weight.
	if summary.GoalMl != 2450 {
		t.Errorf("GoalMl = %d, want 2450", summary.GoalMl)
	}
}

func TestGetDaily_WeightLookupFailureStillServesTheSummary(t *testing.T) {
	repo := mock.NewEntryRepository()
	repo.CurrentWeightKgErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	summary, err := svc.GetDaily(context.Background(), testProfileID, "2026-04-01")
	if err != nil {
		t.Fatalf("GetDaily() error = %v, want nil", err)
	}

	// The previous implementation silently fell back to the default weight.
	if summary.GoalMl != 2450 {
		t.Errorf("GoalMl = %d, want 2450", summary.GoalMl)
	}
}

func TestGetDaily_InvalidDate(t *testing.T) {
	svc := newTestService(t, mock.NewEntryRepository())

	_, err := svc.GetDaily(context.Background(), testProfileID, "01/04/2026")

	if appErr := requireAppError(t, err); appErr.Status != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusBadRequest)
	}
	if appErr := requireAppError(t, err); appErr.Message != "Invalid date format" {
		t.Errorf("message = %q, want %q", appErr.Message, "Invalid date format")
	}
}

func TestGetDaily_MissingDateMeansToday(t *testing.T) {
	svc := newTestService(t, mock.NewEntryRepository())

	summary, err := svc.GetDaily(context.Background(), testProfileID, "")
	if err != nil {
		t.Fatalf("GetDaily() error = %v, want nil", err)
	}

	if want := time.Now().Format("2006-01-02"); summary.Date != want {
		t.Errorf("Date = %q, want today (%q)", summary.Date, want)
	}
}

func TestGetDaily_EmptyEntriesStayNull(t *testing.T) {
	svc := newTestService(t, mock.NewEntryRepository())

	summary, err := svc.GetDaily(context.Background(), testProfileID, "2026-04-01")
	if err != nil {
		t.Fatalf("GetDaily() error = %v, want nil", err)
	}

	if summary.TotalMl != 0 {
		t.Errorf("TotalMl = %d, want 0", summary.TotalMl)
	}
	// A nil slice is what the API has always serialised as JSON null.
	if summary.Entries != nil {
		t.Errorf("Entries = %#v, want nil", summary.Entries)
	}
}

func TestGetDaily_EntriesFailure(t *testing.T) {
	repo := mock.NewEntryRepository()
	repo.GetByDayErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetDaily(context.Background(), testProfileID, "2026-04-01")

	requireStatus(t, err, http.StatusInternalServerError)
}
