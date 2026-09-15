// Package mock holds hand-written test doubles for the user slice, filling the
// role mockery-generated mocks play in the reference architecture.
package mock

import (
	"context"
	"time"

	"fittracker-api/internal/user"
)

// UserRepository is an in-memory user.Repository. Set any of the Err fields to
// force that call to fail; counters record what the service did.
type UserRepository struct {
	SaveErr            error
	GetByEmailErr      error
	GetByIDErr         error
	GetPasswordHashErr error
	UpdateLastLoginErr error
	UpdatePasswordErr  error

	LastLoginUpdates int
	PasswordUpdates  int
	LastSavedHash    string

	accounts map[string]*user.User
	hashes   map[string]string
}

var _ user.Repository = (*UserRepository)(nil)

// NewUserRepository builds an empty in-memory repository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		accounts: map[string]*user.User{},
		hashes:   map[string]string{},
	}
}

// Seed registers an account that appears to already exist.
func (r *UserRepository) Seed(u *user.User, passwordHash string) {
	r.accounts[u.ID] = u
	r.hashes[u.ID] = passwordHash
}

// Hash returns the stored hash of an account, or "" when unknown.
func (r *UserRepository) Hash(id string) string { return r.hashes[id] }

// Save stores the account and its hash, rejecting a duplicate email.
func (r *UserRepository) Save(_ context.Context, u *user.User, passwordHash string) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}

	for _, existing := range r.accounts {
		if existing.Email == u.Email {
			return user.ErrRegistrationFailed
		}
	}

	r.accounts[u.ID] = u
	r.hashes[u.ID] = passwordHash
	r.LastSavedHash = passwordHash
	return nil
}

// GetByEmail returns the account with the given email. Mirroring the SQL that
// selects password_hash, the returned account carries the stored hash.
func (r *UserRepository) GetByEmail(_ context.Context, email string) (*user.User, bool, error) {
	if r.GetByEmailErr != nil {
		return nil, false, r.GetByEmailErr
	}

	for _, account := range r.accounts {
		if account.Email == email {
			found := *account
			found.PasswordHash = r.hashes[found.ID]
			return &found, true, nil
		}
	}
	return nil, false, nil
}

// GetByID returns the account with the given id.
func (r *UserRepository) GetByID(_ context.Context, id string) (*user.User, bool, error) {
	if r.GetByIDErr != nil {
		return nil, false, r.GetByIDErr
	}

	account, ok := r.accounts[id]
	if !ok {
		return nil, false, nil
	}
	return account, true, nil
}

// GetPasswordHash returns the stored hash for an account.
func (r *UserRepository) GetPasswordHash(_ context.Context, id string) (string, bool, error) {
	if r.GetPasswordHashErr != nil {
		return "", false, r.GetPasswordHashErr
	}

	hash, ok := r.hashes[id]
	if !ok {
		return "", false, nil
	}
	return hash, true, nil
}

// UpdateLastLogin counts the stamp.
func (r *UserRepository) UpdateLastLogin(_ context.Context, _ string, _ time.Time) error {
	if r.UpdateLastLoginErr != nil {
		return r.UpdateLastLoginErr
	}
	r.LastLoginUpdates++
	return nil
}

// UpdatePassword replaces the stored hash.
func (r *UserRepository) UpdatePassword(_ context.Context, id, passwordHash string) error {
	if r.UpdatePasswordErr != nil {
		return r.UpdatePasswordErr
	}
	r.hashes[id] = passwordHash
	r.PasswordUpdates++
	return nil
}
