// Package infrastructure implements the user slice outbound ports on SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"fittracker-api/internal/user"
)

// UserRepository is the SQLite implementation of user.Repository.
type UserRepository struct {
	db *sql.DB
}

var _ user.Repository = (*UserRepository)(nil)

// NewUserRepository builds the account repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save inserts the profile and the account in a single transaction. The
// error message for a failing account insert stays generic so callers cannot
// enumerate registered emails.
func (r *UserRepository) Save(ctx context.Context, u *user.User, passwordHash string) error {
	profileID := ""
	if u.ProfileID != nil {
		profileID = *u.ProfileID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO UserProfiles (id, name) VALUES (?, ?)`, profileID, u.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO Users (id, email, password_hash, name, user_profile_id) VALUES (?, ?, ?, ?, ?)`,
		u.ID,
		u.Email,
		passwordHash,
		u.Name,
		profileID,
	)
	if err != nil {
		tx.Rollback()
		return user.ErrRegistrationFailed
	}

	return tx.Commit()
}

// GetByEmail loads an account by email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, bool, error) {
	var u user.User
	var profileID sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password_hash, name, user_profile_id FROM Users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &profileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if profileID.Valid {
		u.ProfileID = &profileID.String
	}

	return &u, true, nil
}

// GetByID loads an account by id.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, bool, error) {
	var u user.User
	var profileID sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, name, user_profile_id FROM Users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &profileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if profileID.Valid {
		u.ProfileID = &profileID.String
	}

	return &u, true, nil
}

// GetPasswordHash returns the stored bcrypt hash for an account.
func (r *UserRepository) GetPasswordHash(ctx context.Context, id string) (string, bool, error) {
	var hash string

	err := r.db.QueryRowContext(ctx, `SELECT password_hash FROM Users WHERE id = ?`, id).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}

	return hash, true, nil
}

// UpdateLastLogin stamps the last successful login.
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE Users SET last_login_at = ? WHERE id = ?`, at, id)
	return err
}

// UpdatePassword replaces the stored password hash.
func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE Users SET password_hash = ? WHERE id = ?`, passwordHash, id)
	return err
}
