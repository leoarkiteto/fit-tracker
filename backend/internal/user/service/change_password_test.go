package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/user/mock"
	"fittracker-api/pkg/auth"
)

func TestChangePassword_Success(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	err := svc.ChangePassword(context.Background(), "user-1", "secret123", "new-secret456")
	if err != nil {
		t.Fatalf("ChangePassword() error = %v, want nil", err)
	}

	if repo.PasswordUpdates != 1 {
		t.Errorf("PasswordUpdates = %d, want 1", repo.PasswordUpdates)
	}
	if !auth.CheckPasswordHash("new-secret456", repo.Hash("user-1")) {
		t.Error("stored hash does not verify against the new password")
	}
	if auth.CheckPasswordHash("secret123", repo.Hash("user-1")) {
		t.Error("stored hash still verifies against the old password")
	}
}

func TestChangePassword_EmptyNewPassword(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	err := svc.ChangePassword(context.Background(), "user-1", "secret123", "")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusBadRequest)
	}
	if appErr.Message != "New password is required" {
		t.Errorf("message = %q, want %q", appErr.Message, "New password is required")
	}
	if repo.PasswordUpdates != 0 {
		t.Errorf("PasswordUpdates = %d, want 0", repo.PasswordUpdates)
	}
}

func TestChangePassword_UnknownUser(t *testing.T) {
	repo := mock.NewUserRepository()
	svc, _ := newTestService(t, repo)

	err := svc.ChangePassword(context.Background(), "missing", "secret123", "new-secret456")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusNotFound)
	}
	if appErr.Message != "User not found" {
		t.Errorf("message = %q, want %q", appErr.Message, "User not found")
	}
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	err := svc.ChangePassword(context.Background(), "user-1", "wrong-password", "new-secret456")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusUnauthorized)
	}
	if appErr.Message != "Current password is incorrect" {
		t.Errorf("message = %q, want %q", appErr.Message, "Current password is incorrect")
	}
	if repo.PasswordUpdates != 0 {
		t.Errorf("PasswordUpdates = %d, want 0", repo.PasswordUpdates)
	}
}

func TestChangePassword_RepositoryFailure(t *testing.T) {
	repo := seededRepository(t, "secret123")
	repo.UpdatePasswordErr = errors.New("database is gone")
	svc, _ := newTestService(t, repo)

	err := svc.ChangePassword(context.Background(), "user-1", "secret123", "new-secret456")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
