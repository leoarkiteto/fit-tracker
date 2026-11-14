package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/water/mock"
)

func TestDelete_ScopesToProfile(t *testing.T) {
	repo := mock.NewEntryRepository()
	svc := newTestService(t, repo)

	if err := svc.Delete(context.Background(), "entry-1", testProfileID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	if repo.DeletedID != "entry-1" || repo.DeletedFor != testProfileID {
		t.Errorf(
			"repository got (id=%q, profile=%q), want (%q, %q)",
			repo.DeletedID,
			repo.DeletedFor,
			"entry-1",
			testProfileID,
		)
	}
}

func TestDelete_UnknownEntryIsNotAnError(t *testing.T) {
	svc := newTestService(t, mock.NewEntryRepository())

	// The previous implementation answered 204 for an unknown id and this keeps
	// that behaviour.
	if err := svc.Delete(context.Background(), "missing", testProfileID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
}

func TestDelete_RepositoryFailure(t *testing.T) {
	repo := mock.NewEntryRepository()
	repo.DeleteErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	err := svc.Delete(context.Background(), "entry-1", testProfileID)

	requireStatus(t, err, http.StatusInternalServerError)
}
