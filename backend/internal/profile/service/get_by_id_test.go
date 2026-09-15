package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/profile/mock"
)

func TestGetByID_Success(t *testing.T) {
	repo, seeded := seededRepository(t)
	svc := newTestService(t, repo)

	result, err := svc.GetByID(context.Background(), seeded.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}

	if result.Name != "Leo" {
		t.Errorf("Name = %q, want %q", result.Name, "Leo")
	}
	if result.Age == nil || *result.Age != 31 {
		t.Errorf("Age = %v, want 31", result.Age)
	}
	if result.CurrentWeight == nil || *result.CurrentWeight != 78.5 {
		t.Errorf("CurrentWeight = %v, want 78.5", result.CurrentWeight)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := newTestService(t, mock.NewProfileRepository())

	_, err := svc.GetByID(context.Background(), "missing")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr.Message != "Profile not found" {
		t.Errorf("message = %q, want %q", appErr.Message, "Profile not found")
	}
}

func TestGetByID_RepositoryFailure(t *testing.T) {
	repo, _ := seededRepository(t)
	repo.GetByIDErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetByID(context.Background(), "profile-1")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
