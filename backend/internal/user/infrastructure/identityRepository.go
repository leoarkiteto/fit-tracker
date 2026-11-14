package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"fittracker-api/internal/user"
)

// IdentityRepository is the SQLite implementation of user.IdentityRepository.
type IdentityRepository struct {
	db *sql.DB
}

var _ user.IdentityRepository = (*IdentityRepository)(nil)

// NewIdentityRepository builds the linked identity repository.
func NewIdentityRepository(db *sql.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

// FindUserByIdentity loads the account linked to (provider, subject).
func (r *IdentityRepository) FindUserByIdentity(ctx context.Context, provider, subject string) (*user.User, bool, error) {
	var u user.User
	var profileID sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		`SELECT u.id, u.email, u.password_hash, u.name, u.user_profile_id
		FROM UserIdentities i JOIN Users u ON u.id = i.user_id
		WHERE i.provider = ? AND i.subject = ?`,
		provider, subject,
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

// FindUsersByEmailFold returns every account whose email equals email, ignoring case.
func (r *IdentityRepository) FindUsersByEmailFold(ctx context.Context, email string) ([]user.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, email, password_hash, name, user_profile_id FROM Users WHERE lower(email) = lower(?)`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		var profileID sql.NullString
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &profileID); err != nil {
			return nil, err
		}
		if profileID.Valid {
			p := profileID.String
			u.ProfileID = &p
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// CreateWithIdentity inserts profile, password-less account and identity in one transaction.
func (r *IdentityRepository) CreateWithIdentity(ctx context.Context, u *user.User, identity user.LinkedIdentity) error {
	profileID := ""
	if u.ProfileID != nil {
		profileID = *u.ProfileID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO UserProfiles (id, name) VALUES (?, ?)`, profileID, u.Name); err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO Users (id, email, password_hash, name, user_profile_id) VALUES (?, ?, ?, ?, ?)`,
		u.ID, u.Email, user.NoPasswordHash, u.Name, profileID,
	); err != nil {
		return err
	}
	if err := insertIdentity(ctx, tx, identity); err != nil {
		return err
	}
	return tx.Commit()
}

// LinkIdentity atomically inserts identity, lowercases the account email and,
// if the account has a password, clears it and deletes its reset codes.
func (r *IdentityRepository) LinkIdentity(ctx context.Context, identity user.LinkedIdentity) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertIdentity(ctx, tx, identity); err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE Users SET email = lower(email) WHERE id = ?`, identity.UserID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return fmt.Errorf("link identity: account %q not found", identity.UserID)
	}

	res, err = tx.ExecContext(
		ctx,
		`UPDATE Users SET password_hash = ? WHERE id = ? AND password_hash <> ?`,
		user.NoPasswordHash, identity.UserID, user.NoPasswordHash,
	)
	if err != nil {
		return err
	}
	changed, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if changed > 0 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM PasswordResetCodes WHERE user_id = ?`, identity.UserID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func insertIdentity(ctx context.Context, tx *sql.Tx, identity user.LinkedIdentity) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO UserIdentities (provider, subject, user_id, created_at) VALUES (?, ?, ?, ?)`,
		identity.Provider, identity.Subject, identity.UserID, identity.CreatedAt.UTC(),
	)
	return err
}
