package infrastructure

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/database"
)

func identityOpenDB(t *testing.T) *sql.DB {
	t.Helper()
	db := database.InitDB(filepath.Join(t.TempDir(), "test.db"))
	t.Cleanup(func() { db.Close() })
	return db
}

// identitySeed saves an account; an empty hash makes it password-less.
func identitySeed(t *testing.T, db *sql.DB, id, email, hash string) {
	t.Helper()
	profileID := "profile-" + id
	err := NewUserRepository(db).Save(context.Background(), &user.User{
		ID: id, Email: email, Name: "N " + id, ProfileID: &profileID,
	}, hash)
	if err != nil {
		t.Fatalf("seed %s: %v", id, err)
	}
}

func identityLink(userID, provider, subject string) user.LinkedIdentity {
	return user.LinkedIdentity{UserID: userID, Provider: provider, Subject: subject, CreatedAt: time.Now()}
}

func identityCount(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func identitySeedCode(t *testing.T, db *sql.DB, userID string) {
	t.Helper()
	c := user.PasswordResetCode{UserID: userID, CodeHash: "h", CreatedAt: time.Now()}
	if err := NewPasswordResetRepository(db).Replace(context.Background(), c); err != nil {
		t.Fatal(err)
	}
}

func TestIdentityDuplicateRejectedAndMultipleProviders(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "u1", "a@b.com", "")
	identitySeed(t, db, "u2", "c@d.com", "")

	if err := repo.LinkIdentity(ctx, identityLink("u1", user.ProviderGoogle, "sub1")); err != nil {
		t.Fatal(err)
	}
	if err := repo.LinkIdentity(ctx, identityLink("u1", user.ProviderGitHub, "42")); err != nil {
		t.Fatalf("second provider on same account: %v", err)
	}
	if err := repo.LinkIdentity(ctx, identityLink("u2", user.ProviderGoogle, "sub1")); err == nil {
		t.Fatal("duplicate (provider, subject) for another account must fail")
	}
	got, found, err := repo.FindUserByIdentity(ctx, user.ProviderGoogle, "sub1")
	if err != nil || !found || got.ID != "u1" {
		t.Fatalf("got (%+v, %v, %v), want u1", got, found, err)
	}
}

func TestIdentityLegacyDatabaseGainsTable(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "old.db")
	old, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = old.Exec(`
		CREATE TABLE UserProfiles (id TEXT PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE Users (
			id TEXT PRIMARY KEY, email TEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL,
			name TEXT NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_login_at DATETIME, user_profile_id TEXT);
		INSERT INTO UserProfiles (id, name) VALUES ('p1', 'Old');
		INSERT INTO Users (id, email, password_hash, name, user_profile_id) VALUES ('u1', 'old@x.com', 'hash', 'Old', 'p1');`)
	if err != nil {
		t.Fatal(err)
	}
	old.Close()

	db := database.InitDB(path)
	defer db.Close()

	if n := identityCount(t, db, `SELECT count(*) FROM sqlite_master WHERE name = 'UserIdentities'`); n != 1 {
		t.Fatalf("UserIdentities missing after InitDB")
	}
	u, found, err := NewUserRepository(db).GetByEmail(ctx, "old@x.com")
	if err != nil || !found || u.PasswordHash != "hash" {
		t.Fatalf("old user: (%+v, %v, %v)", u, found, err)
	}
}

func TestIdentityLinkClearsPasswordAndCodes(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "user-1", "a@b.com", "hash")
	identitySeedCode(t, db, "user-1")

	if err := repo.LinkIdentity(ctx, identityLink("user-1", user.ProviderGoogle, "s")); err != nil {
		t.Fatal(err)
	}
	u, found, err := repo.FindUserByIdentity(ctx, user.ProviderGoogle, "s")
	if err != nil || !found || !user.IsPasswordless(u.PasswordHash) {
		t.Fatalf("got (%+v, %v, %v), want passwordless", u, found, err)
	}
	if _, found, _ := NewPasswordResetRepository(db).Get(ctx, "user-1"); found {
		t.Fatal("reset code must be deleted")
	}
}

func TestIdentityLinkPasswordlessKeepsCodes(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "user-1", "a@b.com", "")
	identitySeedCode(t, db, "user-1")

	if err := repo.LinkIdentity(ctx, identityLink("user-1", user.ProviderGitHub, "7")); err != nil {
		t.Fatal(err)
	}
	if _, found, _ := NewPasswordResetRepository(db).Get(ctx, "user-1"); !found {
		t.Fatal("reset code of a password-less account must survive")
	}
}

func TestIdentityLinkFailureRollsBackEverything(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "user-1", "Mixed@B.com", "hash")
	identitySeed(t, db, "user-2", "o@b.com", "")
	identitySeedCode(t, db, "user-1")

	// Unknown account: the identity insert may pass but the email update hits 0 rows.
	if err := repo.LinkIdentity(ctx, identityLink("ghost", user.ProviderGoogle, "s")); err == nil {
		t.Fatal("expected error for unknown account")
	}
	// Duplicate identity owned by another account: insert fails.
	if err := repo.LinkIdentity(ctx, identityLink("user-2", user.ProviderGoogle, "dup")); err != nil {
		t.Fatal(err)
	}
	if err := repo.LinkIdentity(ctx, identityLink("user-1", user.ProviderGoogle, "dup")); err == nil {
		t.Fatal("expected duplicate error")
	}

	if n := identityCount(t, db, `SELECT count(*) FROM UserIdentities WHERE user_id IN ('user-1', 'ghost')`); n != 0 {
		t.Fatalf("identity rows leaked: %d", n)
	}
	var email, hash string
	if err := db.QueryRow(`SELECT email, password_hash FROM Users WHERE id = 'user-1'`).Scan(&email, &hash); err != nil {
		t.Fatal(err)
	}
	if email != "Mixed@B.com" || hash != "hash" {
		t.Fatalf("account changed: %q %q", email, hash)
	}
	if _, found, _ := NewPasswordResetRepository(db).Get(ctx, "user-1"); !found {
		t.Fatal("reset code must be untouched")
	}
}

// AC-24: a failure in the last step of the link (deleting the reset codes),
// after the identity insert, the email lowercase and the password clear have
// run inside the transaction, must roll all of them back.
func TestIdentityLinkFailureAfterPasswordClearRollsBack(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "user-7", "Mixed@B.com", "hash")
	identitySeedCode(t, db, "user-7")

	if _, err := db.Exec(`CREATE TRIGGER identity_test_fail_code_delete
		BEFORE DELETE ON PasswordResetCodes
		BEGIN SELECT RAISE(ABORT, 'forced failure'); END;`); err != nil {
		t.Fatal(err)
	}

	if err := repo.LinkIdentity(ctx, identityLink("user-7", user.ProviderGoogle, "s7")); err == nil {
		t.Fatal("expected the forced failure to surface")
	}

	if n := identityCount(t, db, `SELECT count(*) FROM UserIdentities WHERE subject = 's7'`); n != 0 {
		t.Fatalf("identity row leaked: %d", n)
	}
	var email, hash string
	if err := db.QueryRow(`SELECT email, password_hash FROM Users WHERE id = 'user-7'`).Scan(&email, &hash); err != nil {
		t.Fatal(err)
	}
	if email != "Mixed@B.com" || hash != "hash" {
		t.Fatalf("account changed despite rollback: email %q hash %q", email, hash)
	}
	if _, found, _ := NewPasswordResetRepository(db).Get(ctx, "user-7"); !found {
		t.Fatal("reset code must be untouched")
	}
}

func TestIdentityLinkLowercasesEmail(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "user-1", "Vitima@Gmail.com", "hash")

	if err := repo.LinkIdentity(ctx, identityLink("user-1", user.ProviderGoogle, "s")); err != nil {
		t.Fatal(err)
	}
	u, found, err := NewUserRepository(db).GetByEmail(ctx, "vitima@gmail.com")
	if err != nil || !found || u.ID != "user-1" {
		t.Fatalf("got (%+v, %v, %v)", u, found, err)
	}
}

func TestIdentityFindUsersByEmailFold(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "u1", "Ana@X.com", "hash")
	identitySeed(t, db, "u2", "ana@x.com", "")
	identitySeed(t, db, "u3", "other@x.com", "")

	got, err := repo.FindUsersByEmailFold(ctx, "ANA@x.COM")
	if err != nil || len(got) != 2 {
		t.Fatalf("got (%+v, %v), want 2 matches", got, err)
	}
	for _, u := range got {
		if u.ProfileID == nil || (u.ID == "u1" && u.PasswordHash != "hash") {
			t.Fatalf("incomplete account: %+v", u)
		}
	}
	got, err = repo.FindUsersByEmailFold(ctx, "none@x.com")
	if err != nil || len(got) != 0 {
		t.Fatalf("got (%+v, %v), want none", got, err)
	}
}

func TestIdentityFindUserByIdentityNotFound(t *testing.T) {
	db := identityOpenDB(t)
	got, found, err := NewIdentityRepository(db).FindUserByIdentity(context.Background(), user.ProviderGoogle, "x")
	if err != nil || found || got != nil {
		t.Fatalf("got (%v, %v, %v), want (nil, false, nil)", got, found, err)
	}
}

func TestIdentityCreateWithIdentity(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	pid := "p-new"
	u := &user.User{ID: "new", Email: "new@x.com", Name: "New", ProfileID: &pid}

	if err := repo.CreateWithIdentity(ctx, u, identityLink("new", user.ProviderGoogle, "s")); err != nil {
		t.Fatal(err)
	}
	got, found, err := repo.FindUserByIdentity(ctx, user.ProviderGoogle, "s")
	if err != nil || !found || got.ProfileID == nil || *got.ProfileID != pid || !user.IsPasswordless(got.PasswordHash) {
		t.Fatalf("got (%+v, %v, %v)", got, found, err)
	}
}

func TestIdentityCreateWithIdentityRollsBackOnDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	db := identityOpenDB(t)
	repo := NewIdentityRepository(db)
	identitySeed(t, db, "u1", "dup@x.com", "hash")
	pid := "p-new"
	u := &user.User{ID: "new", Email: "dup@x.com", Name: "New", ProfileID: &pid}

	if err := repo.CreateWithIdentity(ctx, u, identityLink("new", user.ProviderGoogle, "s")); err == nil {
		t.Fatal("expected error for existing email")
	}
	if n := identityCount(t, db, `SELECT count(*) FROM UserProfiles WHERE id = 'p-new'`); n != 0 {
		t.Fatalf("profile left behind")
	}
	if n := identityCount(t, db, `SELECT count(*) FROM UserIdentities`); n != 0 {
		t.Fatalf("identity left behind")
	}
}
