package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/user/mock"
)

func TestMe_Success(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	result, err := svc.Me(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Me() error = %v, want nil", err)
	}

	if result.ID != "user-1" {
		t.Errorf("ID = %q, want %q", result.ID, "user-1")
	}
	if result.Email != "leo@example.com" {
		t.Errorf("Email = %q, want %q", result.Email, "leo@example.com")
	}
	if result.ProfileID == nil || *result.ProfileID != "profile-1" {
		t.Errorf("ProfileID = %v, want %q", result.ProfileID, "profile-1")
	}
}

func TestMe_NotFound(t *testing.T) {
	repo := mock.NewUserRepository()
	svc, _ := newTestService(t, repo)

	_, err := svc.Me(context.Background(), "missing")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr.Message != "User not found" {
		t.Errorf("message = %q, want %q", appErr.Message, "User not found")
	}
}

func TestMe_RepositoryFailure(t *testing.T) {
	repo := seededRepository(t, "secret123")
	repo.GetByIDErr = errors.New("database is gone")
	svc, _ := newTestService(t, repo)

	_, err := svc.Me(context.Background(), "user-1")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
