package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/internal/bioimpedance/mock"
)

func TestCreate_DefaultsDateAndFillsIdentifiers(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	svc := newTestService(t, repo)

	before := time.Now()
	measurement := bioimpedance.Measurement{Weight: 78.2, BodyFatPercentage: 18.4}

	result, err := svc.Create(context.Background(), testProfileID, &measurement)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if repo.SaveCalls != 1 {
		t.Errorf("SaveCalls = %d, want 1", repo.SaveCalls)
	}
	if result.ID == "" {
		t.Error("ID = empty, want a generated id")
	}
	if result.UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", result.UserProfileID, testProfileID)
	}
	if result.Date.Before(before) || result.Date.After(time.Now()) {
		t.Errorf("Date = %v, want the moment of the call", result.Date)
	}
	if repo.Saved == nil || repo.Saved.ID != result.ID {
		t.Error("the repository did not receive the created measurement")
	}
}

func TestCreate_KeepsProvidedDate(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	svc := newTestService(t, repo)

	date := time.Date(2026, 4, 1, 7, 30, 0, 0, time.UTC)
	measurement := bioimpedance.Measurement{Date: date, Weight: 78.2}

	result, err := svc.Create(context.Background(), testProfileID, &measurement)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if !result.Date.Equal(date) {
		t.Errorf("Date = %v, want the provided %v", result.Date, date)
	}
}

func TestCreate_RepositoryFailure(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	repo.SaveErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	measurement := bioimpedance.Measurement{Weight: 78.2}

	_, err := svc.Create(context.Background(), testProfileID, &measurement)

	requireStatus(t, err, http.StatusInternalServerError)
}
