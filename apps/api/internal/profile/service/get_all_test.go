package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/profile/mock"
)

func TestGetAll_ReturnsProfiles(t *testing.T) {
	repo, _ := seededRepository(t)
	svc := newTestService(t, repo)

	profiles, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	if len(profiles) != 1 {
		t.Fatalf("len(profiles) = %d, want 1", len(profiles))
	}
	if profiles[0].ID != "profile-1" {
		t.Errorf("profiles[0].ID = %q, want %q", profiles[0].ID, "profile-1")
	}
	if profiles[0].ExperienceLevel != "Beginner" {
		t.Errorf("ExperienceLevel = %q, want %q", profiles[0].ExperienceLevel, "Beginner")
	}
}

func TestGetAll_EmptyStaysNull(t *testing.T) {
	svc := newTestService(t, mock.NewProfileRepository())

	profiles, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}

	// A nil slice is what the API has always serialised as JSON null.
	if profiles != nil {
		t.Errorf("profiles = %#v, want nil", profiles)
	}
}

func TestGetAll_RepositoryFailure(t *testing.T) {
	repo, _ := seededRepository(t)
	repo.GetAllErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.GetAll(context.Background())

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
