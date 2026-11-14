package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestDelete_RemovesProfile(t *testing.T) {
	repo, seeded := seededRepository(t)
	svc := newTestService(t, repo)

	if err := svc.Delete(context.Background(), seeded.ID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	if repo.DeleteCalls != 1 {
		t.Errorf("DeleteCalls = %d, want 1", repo.DeleteCalls)
	}

	_, found, err := repo.GetByID(context.Background(), seeded.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}
	if found {
		t.Error("profile still present after Delete()")
	}
}

func TestDelete_UnknownProfileIsNotAnError(t *testing.T) {
	repo, _ := seededRepository(t)
	svc := newTestService(t, repo)

	// The previous implementation answered 204 for an unknown id and this keeps
	// that behaviour.
	if err := svc.Delete(context.Background(), "missing"); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
}

func TestDelete_RepositoryFailure(t *testing.T) {
	repo, seeded := seededRepository(t)
	repo.DeleteErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	err := svc.Delete(context.Background(), seeded.ID)

	appErr := requireAppError(t, err)
	if appErr.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", appErr.Status, http.StatusInternalServerError)
	}
}
