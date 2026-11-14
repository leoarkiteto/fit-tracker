package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestUpdate_StoresFields(t *testing.T) {
	repo, seeded := seededRepository(t)
	svc := newTestService(t, repo)

	goal := 72.0
	updated := *seeded
	updated.Name = "Leonardo"
	updated.GoalWeight = &goal

	if err := svc.Update(context.Background(), &updated); err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	if repo.UpdateCalls != 1 {
		t.Errorf("UpdateCalls = %d, want 1", repo.UpdateCalls)
	}

	stored, found, err := repo.GetByID(context.Background(), seeded.ID)
	if err != nil || !found {
		t.Fatalf("GetByID() = (%v, %v), want the stored profile", found, err)
	}
	if stored.Name != "Leonardo" {
		t.Errorf("Name = %q, want %q", stored.Name, "Leonardo")
	}
	if stored.GoalWeight == nil || *stored.GoalWeight != 72.0 {
		t.Errorf("GoalWeight = %v, want 72", stored.GoalWeight)
	}
}

func TestUpdate_RepositoryFailure(t *testing.T) {
	repo, seeded := seededRepository(t)
	repo.UpdateErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	err := svc.Update(context.Background(), seeded)

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
