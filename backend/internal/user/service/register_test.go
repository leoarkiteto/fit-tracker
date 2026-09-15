package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/mock"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/auth"
)

func TestRegister_Success(t *testing.T) {
	repo := mock.NewUserRepository()
	svc, signer := newTestService(t, repo)

	result, err := svc.Register(context.Background(), "leo@example.com", "secret123", "Leo")
	if err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	if result.User.Email != "leo@example.com" {
		t.Errorf("User.Email = %q, want %q", result.User.Email, "leo@example.com")
	}
	if result.User.Name != "Leo" {
		t.Errorf("User.Name = %q, want %q", result.User.Name, "Leo")
	}
	if result.User.ProfileID == nil || *result.User.ProfileID == "" {
		t.Fatal("User.ProfileID = empty, want a generated profile id")
	}
	if result.Token == "" {
		t.Error("Token = empty, want a signed token")
	}
	if !result.ExpiresAt.After(time.Now()) {
		t.Errorf("ExpiresAt = %v, want a future expiry", result.ExpiresAt)
	}

	// The password must be stored hashed, never in plaintext.
	if repo.LastSavedHash == "secret123" {
		t.Error("password was stored in plaintext")
	}
	if !auth.CheckPasswordHash("secret123", repo.LastSavedHash) {
		t.Error("stored hash does not verify against the plaintext password")
	}

	claims, err := signer.ValidateToken(result.Token)
	if err != nil {
		t.Fatalf("issued token does not validate: %v", err)
	}
	if got := (*claims)["sub"]; got != result.User.ID {
		t.Errorf("token sub = %v, want %v", got, result.User.ID)
	}
	if got := (*claims)["profileId"]; got != *result.User.ProfileID {
		t.Errorf("token profileId = %v, want %v", got, *result.User.ProfileID)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := mock.NewUserRepository()
	repo.Seed(&user.User{ID: "existing-user", Email: "leo@example.com"}, "existing-hash")
	svc, _ := newTestService(t, repo)

	_, err := svc.Register(context.Background(), "leo@example.com", "secret123", "Leo")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusConflict {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusConflict)
	}
	// The message must stay generic so emails cannot be enumerated.
	if appErr.Message != "Registration failed" {
		t.Errorf("message = %q, want %q", appErr.Message, "Registration failed")
	}
}

func TestRegister_RepositoryFailure(t *testing.T) {
	repo := mock.NewUserRepository()
	repo.SaveErr = errors.New("database is gone")
	svc, _ := newTestService(t, repo)

	_, err := svc.Register(context.Background(), "leo@example.com", "secret123", "Leo")

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
	if appErr.Message != "Internal server error" {
		t.Errorf("message = %q, want %q", appErr.Message, "Internal server error")
	}
}

// requireAppError asserts that err carries an HTTP status and returns it.
func requireAppError(t *testing.T, err error) *apperror.Error {
	t.Helper()

	if err == nil {
		t.Fatal("error = nil, want an *apperror.Error")
	}

	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v (%T), want *apperror.Error", err, err)
	}

	return appErr
}
