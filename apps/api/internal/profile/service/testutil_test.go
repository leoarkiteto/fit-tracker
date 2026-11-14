package service

import (
	"errors"
	"testing"

	"fittracker-api/internal/profile"
	"fittracker-api/internal/profile/mock"
	"fittracker-api/pkg/apperror"
)

// newTestService builds the profile service on an in-memory repository.
func newTestService(t *testing.T, repo profile.Repository) *Service {
	t.Helper()
	return New(repo)
}

// seededRepository returns a repository holding one profile.
func seededRepository(t *testing.T) (*mock.ProfileRepository, *profile.Profile) {
	t.Helper()

	age := 31
	weight := 78.5
	seeded := &profile.Profile{
		ID:                   "profile-1",
		Name:                 "Leo",
		Age:                  &age,
		CurrentWeight:        &weight,
		ExperienceLevel:      "Beginner",
		EquipmentType:        "Gym",
		AvailableDaysPerWeek: 3,
	}

	repo := mock.NewProfileRepository()
	repo.Seed(seeded)

	return repo, seeded
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
