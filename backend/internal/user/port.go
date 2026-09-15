package user

import (
	"context"
	"time"
)

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
