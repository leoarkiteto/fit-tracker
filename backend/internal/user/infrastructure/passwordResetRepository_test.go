package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/database"
)

func newResetRepo(t *testing.T) *PasswordResetRepository {
	t.Helper()
	db := database.InitDB(filepath.Join(t.TempDir(), "test.db"))
	t.Cleanup(func() { db.Close() })

	profileID := "profile-1"
	err := NewUserRepository(db).Save(context.Background(), &user.User{
		ID:        "user-1",
		Email:     "a@b.com",
		Name:      "A",
		ProfileID: &profileID,
	}, "hash")
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return NewPasswordResetRepository(db)
}

func code(hash string, at time.Time) user.PasswordResetCode {
	return user.PasswordResetCode{UserID: "user-1", CodeHash: hash, CreatedAt: at}
}

func TestGetAbsent(t *testing.T) {
	repo := newResetRepo(t)
	got, found, err := repo.Get(context.Background(), "user-1")
	if err != nil || found || got != nil {
		t.Fatalf("got (%v, %v, %v), want (nil, false, nil)", got, found, err)
	}
}

func TestReplaceOverwrites(t *testing.T) {
	ctx := context.Background()
	repo := newResetRepo(t)
	now := time.Now().UTC()

	if err := repo.Replace(ctx, code("h1", now)); err != nil {
		t.Fatal(err)
	}
	if ok, _ := repo.ConsumeAttempt(ctx, "user-1", "h1", 5); !ok {
		t.Fatal("attempt should be granted")
	}
	if ok, _ := repo.MarkUsed(ctx, "user-1", "h1", now); !ok {
		t.Fatal("mark used should win")
	}

	if err := repo.Replace(ctx, code("h2", now.Add(time.Minute))); err != nil {
		t.Fatal(err)
	}
	got, found, err := repo.Get(ctx, "user-1")
	if err != nil || !found {
		t.Fatalf("get: %v %v", found, err)
	}
	if got.CodeHash != "h2" || got.Attempts != 0 || got.UsedAt != nil {
		t.Fatalf("unexpected row after replace: %+v", got)
	}
}

func TestConsumeAttemptUpToMax(t *testing.T) {
	ctx := context.Background()
	repo := newResetRepo(t)
	if err := repo.Replace(ctx, code("h1", time.Now())); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		ok, err := repo.ConsumeAttempt(ctx, "user-1", "h1", 5)
		if err != nil || !ok {
			t.Fatalf("attempt %d: got (%v, %v), want granted", i+1, ok, err)
		}
	}
	ok, err := repo.ConsumeAttempt(ctx, "user-1", "h1", 5)
	if err != nil || ok {
		t.Fatalf("6th attempt: got (%v, %v), want refused", ok, err)
	}
	got, _, _ := repo.Get(ctx, "user-1")
	if got.Attempts != 5 {
		t.Fatalf("attempts = %d, want 5", got.Attempts)
	}
}

func TestStaleCodeHashRefused(t *testing.T) {
	ctx := context.Background()
	repo := newResetRepo(t)
	if err := repo.Replace(ctx, code("new", time.Now())); err != nil {
		t.Fatal(err)
	}

	if ok, err := repo.ConsumeAttempt(ctx, "user-1", "old", 5); err != nil || ok {
		t.Fatalf("ConsumeAttempt stale: got (%v, %v)", ok, err)
	}
	if ok, err := repo.MarkUsed(ctx, "user-1", "old", time.Now()); err != nil || ok {
		t.Fatalf("MarkUsed stale: got (%v, %v)", ok, err)
	}
}

func TestMarkUsedWonThenLost(t *testing.T) {
	ctx := context.Background()
	repo := newResetRepo(t)
	if err := repo.Replace(ctx, code("h1", time.Now())); err != nil {
		t.Fatal(err)
	}
	at := time.Now().Add(time.Second)

	if ok, err := repo.MarkUsed(ctx, "user-1", "h1", at); err != nil || !ok {
		t.Fatalf("first MarkUsed: got (%v, %v)", ok, err)
	}
	if ok, err := repo.MarkUsed(ctx, "user-1", "h1", at); err != nil || ok {
		t.Fatalf("second MarkUsed: got (%v, %v)", ok, err)
	}
	if ok, _ := repo.ConsumeAttempt(ctx, "user-1", "h1", 5); ok {
		t.Fatal("used code must refuse attempts")
	}
	got, _, _ := repo.Get(ctx, "user-1")
	if got.UsedAt == nil || !got.UsedAt.Equal(at) || got.UsedAt.Location() != time.UTC {
		t.Fatalf("UsedAt = %v, want %v in UTC", got.UsedAt, at.UTC())
	}
}

func TestCreatedAtRoundTrip(t *testing.T) {
	ctx := context.Background()
	repo := newResetRepo(t)
	zone := time.FixedZone("X", -3*3600)
	at := time.Date(2026, 11, 13, 10, 30, 15, 123456789, zone)

	if err := repo.Replace(ctx, code("h1", at)); err != nil {
		t.Fatal(err)
	}
	got, _, err := repo.Get(ctx, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if !got.CreatedAt.Equal(at) {
		t.Fatalf("CreatedAt = %v, want %v", got.CreatedAt, at)
	}
	if got.CreatedAt.Location() != time.UTC {
		t.Fatalf("CreatedAt location = %v, want UTC", got.CreatedAt.Location())
	}
}
