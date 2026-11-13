package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"fittracker-api/internal/user"
)

// PasswordResetRepository is the SQLite implementation of
// user.PasswordResetRepository. Times are written in UTC and read back as UTC.
type PasswordResetRepository struct {
	db *sql.DB
}

var _ user.PasswordResetRepository = (*PasswordResetRepository)(nil)

// NewPasswordResetRepository builds the reset code repository.
func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Replace stores code as the account's only reset code, resetting attempts and
// used_at.
func (r *PasswordResetRepository) Replace(ctx context.Context, code user.PasswordResetCode) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO PasswordResetCodes (user_id, code_hash, created_at, attempts, used_at)
		VALUES (?, ?, ?, 0, NULL)
		ON CONFLICT(user_id) DO UPDATE SET
			code_hash = excluded.code_hash,
			created_at = excluded.created_at,
			attempts = 0,
			used_at = NULL`,
		code.UserID,
		code.CodeHash,
		code.CreatedAt.UTC(),
	)
	return err
}

// Get loads the account's reset code.
func (r *PasswordResetRepository) Get(ctx context.Context, userID string) (*user.PasswordResetCode, bool, error) {
	var c user.PasswordResetCode
	var usedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`SELECT user_id, code_hash, created_at, attempts, used_at FROM PasswordResetCodes WHERE user_id = ?`,
		userID,
	).Scan(&c.UserID, &c.CodeHash, &c.CreatedAt, &c.Attempts, &usedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	c.CreatedAt = c.CreatedAt.UTC()
	if usedAt.Valid {
		t := usedAt.Time.UTC()
		c.UsedAt = &t
	}

	return &c, true, nil
}

// ConsumeAttempt counts one attempt in a single conditional update.
func (r *PasswordResetRepository) ConsumeAttempt(ctx context.Context, userID, codeHash string, max int) (bool, error) {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE PasswordResetCodes SET attempts = attempts + 1
		WHERE user_id = ? AND code_hash = ? AND used_at IS NULL AND attempts < ?`,
		userID,
		codeHash,
		max,
	)
	if err != nil {
		return false, err
	}
	return rowWasChanged(res)
}

// MarkUsed stamps the code as used if it is still unused.
func (r *PasswordResetRepository) MarkUsed(ctx context.Context, userID, codeHash string, at time.Time) (bool, error) {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE PasswordResetCodes SET used_at = ?
		WHERE user_id = ? AND code_hash = ? AND used_at IS NULL`,
		at.UTC(),
		userID,
		codeHash,
	)
	if err != nil {
		return false, err
	}
	return rowWasChanged(res)
}

func rowWasChanged(res sql.Result) (bool, error) {
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}
