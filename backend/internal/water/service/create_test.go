package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"fittracker-api/internal/water/mock"
)

func TestCreate_FillsIDAndDefaultsConsumedAt(t *testing.T) {
	repo := mock.NewEntryRepository()
	svc := newTestService(t, repo)

	before := time.Now()
	entry, err := svc.Create(context.Background(), testProfileID, 250, nil)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if repo.SaveCalls != 1 {
		t.Errorf("SaveCalls = %d, want 1", repo.SaveCalls)
	}
	if entry.ID == "" {
		t.Error("ID = empty, want a generated id")
	}
	if entry.UserProfileID != testProfileID {
		t.Errorf("UserProfileID = %q, want %q", entry.UserProfileID, testProfileID)
	}
	if entry.AmountMl != 250 {
		t.Errorf("AmountMl = %d, want 250", entry.AmountMl)
	}
	if entry.ConsumedAt.Before(before) || entry.ConsumedAt.After(time.Now()) {
		t.Errorf("ConsumedAt = %v, want the moment of the call", entry.ConsumedAt)
	}
	if repo.Saved == nil || repo.Saved.ID != entry.ID {
		t.Error("the repository did not receive the created entry")
	}
}

func TestCreate_KeepsProvidedConsumedAt(t *testing.T) {
	repo := mock.NewEntryRepository()
	svc := newTestService(t, repo)

	consumedAt := time.Date(2026, 4, 1, 9, 15, 0, 0, time.UTC)

	entry, err := svc.Create(context.Background(), testProfileID, 300, &consumedAt)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if !entry.ConsumedAt.Equal(consumedAt) {
		t.Errorf("ConsumedAt = %v, want the provided %v", entry.ConsumedAt, consumedAt)
	}
}

func TestCreate_RepositoryFailure(t *testing.T) {
	repo := mock.NewEntryRepository()
	repo.SaveErr = errors.New("database is gone")
	svc := newTestService(t, repo)

	_, err := svc.Create(context.Background(), testProfileID, 250, nil)

	requireStatus(t, err, http.StatusInternalServerError)
}
