package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"fittracker-api/internal/bioimpedance/mock"
)

func TestDelete_ScopesToProfile(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	svc := newTestService(t, repo)

	if err := svc.Delete(context.Background(), "measurement-1", testProfileID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	if repo.DeletedID != "measurement-1" || repo.DeletedFor != testProfileID {
		t.Errorf(
			"repository got (id=%q, profile=%q), want (%q, %q)",
			repo.DeletedID,
			repo.DeletedFor,
			"measurement-1",
			testProfileID,
		)
	}
}

func TestDelete_RepositoryFailure(t *testing.T) {
	repo := mock.NewMeasurementRepository()
	repo.DeleteErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	err := svc.Delete(context.Background(), "measurement-1", testProfileID)

	requireStatus(t, err, http.StatusInternalServerError)
}
