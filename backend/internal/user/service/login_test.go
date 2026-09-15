package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/pkg/auth"
)

func seededRepository(t *testing.T, password string) *mock.UserRepository {
	t.Helper()

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("hashing seed password: %v", err)
	}

	profileID := "profile-1"
	repo := mock.NewUserRepository()
	repo.Seed(
		&user.User{ID: "user-1", Email: "leo@example.com", Name: "Leo", ProfileID: &profileID},
		hash,
	)

	return repo
}

func TestLogin_Success(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, signer := newTestService(t, repo)

	result, err := svc.Login(context.Background(), "leo@example.com", "secret123")
	if err != nil {
		t.Fatalf("Login() error = %v, want nil", err)
	}

	if result.User.ID != "user-1" {
		t.Errorf("User.ID = %q, want %q", result.User.ID, "user-1")
	}
	if result.User.ProfileID == nil || *result.User.ProfileID != "profile-1" {
		t.Errorf("User.ProfileID = %v, want %q", result.User.ProfileID, "profile-1")
	}
	if repo.LastLoginUpdates != 1 {
		t.Errorf("LastLoginUpdates = %d, want 1", repo.LastLoginUpdates)
	}

	claims, err := signer.ValidateToken(result.Token)
	if err != nil {
		t.Fatalf("issued token does not validate: %v", err)
	}
	if got := (*claims)["profileId"]; got != "profile-1" {
		t.Errorf("token profileId = %v, want %q", got, "profile-1")
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	_, err := svc.Login(context.Background(), "nobody@example.com", "secret123")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusUnauthorized)
	}
	if appErr.Message != "Invalid credentials" {
		t.Errorf("message = %q, want %q", appErr.Message, "Invalid credentials")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := seededRepository(t, "secret123")
	svc, _ := newTestService(t, repo)

	_, err := svc.Login(context.Background(), "leo@example.com", "wrong-password")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusUnauthorized)
	}
	if repo.LastLoginUpdates != 0 {
		t.Errorf("LastLoginUpdates = %d, want 0", repo.LastLoginUpdates)
	}
}

func TestLogin_FailedLastLoginStampStillSucceeds(t *testing.T) {
	repo := seededRepository(t, "secret123")
	repo.UpdateLastLoginErr = errors.New("write failed")
	svc, _ := newTestService(t, repo)

	if _, err := svc.Login(context.Background(), "leo@example.com", "secret123"); err != nil {
		t.Fatalf("Login() error = %v, want nil", err)
	}
}

func TestLogin_RepositoryFailure(t *testing.T) {
	repo := seededRepository(t, "secret123")
	repo.GetByEmailErr = errors.New("database is gone")
	svc, _ := newTestService(t, repo)

	_, err := svc.Login(context.Background(), "leo@example.com", "secret123")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
