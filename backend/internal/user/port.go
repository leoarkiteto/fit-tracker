package user

import (
	"context"
	"time"
)

// PasswordResetRepository is the outbound port for password reset codes. An
// account has at most one code.
type PasswordResetRepository interface {
	// Replace stores code as the account's only reset code, overwriting any previous one.
	Replace(ctx context.Context, code PasswordResetCode) error
	// Get loads the account's reset code.
	Get(ctx context.Context, userID string) (*PasswordResetCode, bool, error)
	// ConsumeAttempt atomically counts one attempt against the account's code
	// identified by codeHash, if it is unused and has fewer than max attempts.
	// granted reports whether the attempt was counted.
	ConsumeAttempt(ctx context.Context, userID, codeHash string, max int) (granted bool, err error)
	// MarkUsed stamps the code identified by codeHash as used if it is still
	// unused; won reports whether this call did it.
	MarkUsed(ctx context.Context, userID, codeHash string, at time.Time) (won bool, err error)
}

// Mailer is the outbound port that delivers reset codes to the account email.
type Mailer interface {
	SendPasswordResetCode(ctx context.Context, email, code string) error
}

// Repository is the outbound port for account persistence. Implementations live
// in the slice's infrastructure package.
//
// Lookups use the (value, found, error) shape so "not found" stays distinct
// from a genuine storage failure.
type Repository interface {
	// Save inserts the account together with the profile it belongs to.
	// ProfileID on the account must be set.
	Save(ctx context.Context, u *User, passwordHash string) error

	// GetByEmail loads an account by email address.
	GetByEmail(ctx context.Context, email string) (*User, bool, error)

	// GetByID loads an account by id.
	GetByID(ctx context.Context, id string) (*User, bool, error)

	// GetPasswordHash returns the stored hash for an account.
	GetPasswordHash(ctx context.Context, id string) (string, bool, error)

	// UpdateLastLogin stamps the last successful login.
	UpdateLastLogin(ctx context.Context, id string, at time.Time) error

	// UpdatePassword replaces the stored password hash.
	UpdatePassword(ctx context.Context, id, passwordHash string) error
}
