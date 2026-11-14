package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/internal/bioimpedance/mock"
)

func TestGetLatest_Success(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	repo.Latest = &bioimpedance.Measurement{ID: "measurement-2", Weight: 77.4}
	svc := newTestService(t, repo)

	result, err := svc.GetLatest(context.Background(), testProfileID)
	if err != nil {
		t.Fatalf("GetLatest() error = %v, want nil", err)
	}

	if result.ID != "measurement-2" {
		t.Errorf("ID = %q, want %q", result.ID, "measurement-2")
	}
	if result.UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", result.UserProfileID, testProfileID)
	}
}

func TestGetLatest_NoneRecorded(t *testing.T) {
	svc := newTestService(t, mock.NewMeasurementRepository())

	_, err := svc.GetLatest(context.Background(), testProfileID)

	if appErr := requireAppError(t, err); appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr := requireAppError(t, err); appErr.Message != "No data found" {
		t.Errorf("message = %q, want %q", appErr.Message, "No data found")
	}
}

func TestGetLatest_RepositoryFailure(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	repo.GetLatestErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetLatest(context.Background(), testProfileID)

	requireStatus(t, err, http.StatusInternalServerError)
}
