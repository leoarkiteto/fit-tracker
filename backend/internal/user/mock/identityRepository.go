package mock

import (
	"context"
	"strings"

	"fittracker-api/internal/user"
)

// IdentityRepository is an in-memory user.IdentityRepository sharing the state
// of a UserRepository and a PasswordResetRepository, so a link is observable
// through them. Set CreateErr or LinkErr to force that call to fail with no
// side effect (the real adapter is all-or-nothing).
type IdentityRepository struct {
	FindUserByIdentityErr   error
	FindUsersByEmailFoldErr error
	CreateErr               error
	LinkErr                 error

	Creates int
	Links   int

	users      *UserRepository
	resets     *PasswordResetRepository
	identities map[string]user.LinkedIdentity
}

var _ user.IdentityRepository = (*IdentityRepository)(nil)

// NewIdentityRepository builds an empty repository over the given mocks.
func NewIdentityRepository(users *UserRepository, resets *PasswordResetRepository) *IdentityRepository {
	return &IdentityRepository{
		users:      users,
		resets:     resets,
		identities: map[string]user.LinkedIdentity{},
	}
}

func identityKey(provider, subject string) string { return provider + "\x00" + subject }

// SeedIdentity stores a linked identity directly, bypassing LinkIdentity.
func (r *IdentityRepository) SeedIdentity(identity user.LinkedIdentity) {
	r.identities[identityKey(identity.Provider, identity.Subject)] = identity
}

// Identity returns the stored identity for (provider, subject).
func (r *IdentityRepository) Identity(provider, subject string) (user.LinkedIdentity, bool) {
	identity, ok := r.identities[identityKey(provider, subject)]
	return identity, ok
}

// Identities returns every stored identity.
func (r *IdentityRepository) Identities() []user.LinkedIdentity {
	out := make([]user.LinkedIdentity, 0, len(r.identities))
	for _, identity := range r.identities {
		out = append(out, identity)
	}
	return out
}

// FindUserByIdentity returns the linked account, carrying its stored hash.
func (r *IdentityRepository) FindUserByIdentity(_ context.Context, provider, subject string) (*user.User, bool, error) {
	if r.FindUserByIdentityErr != nil {
		return nil, false, r.FindUserByIdentityErr
	}
	identity, ok := r.identities[identityKey(provider, subject)]
	if !ok {
		return nil, false, nil
	}
	account, ok := r.users.accounts[identity.UserID]
	if !ok {
		return nil, false, nil
	}
	found := *account
	found.PasswordHash = r.users.hashes[found.ID]
	return &found, true, nil
}

// FindUsersByEmailFold returns the accounts whose email matches ignoring case.
func (r *IdentityRepository) FindUsersByEmailFold(_ context.Context, email string) ([]user.User, error) {
	if r.FindUsersByEmailFoldErr != nil {
		return nil, r.FindUsersByEmailFoldErr
	}
	var found []user.User
	for _, account := range r.users.accounts {
		if strings.EqualFold(account.Email, email) {
			match := *account
			match.PasswordHash = r.users.hashes[match.ID]
			found = append(found, match)
		}
	}
	return found, nil
}

// CreateWithIdentity stores a password-less account and its identity.
func (r *IdentityRepository) CreateWithIdentity(_ context.Context, u *user.User, identity user.LinkedIdentity) error {
	if r.CreateErr != nil {
		return r.CreateErr
	}
	if _, dup := r.identities[identityKey(identity.Provider, identity.Subject)]; dup {
		return user.ErrRegistrationFailed
	}
	for _, existing := range r.users.accounts {
		if existing.Email == u.Email {
			return user.ErrRegistrationFailed
		}
	}

	r.users.accounts[u.ID] = u
	r.users.hashes[u.ID] = user.NoPasswordHash
	r.identities[identityKey(identity.Provider, identity.Subject)] = identity
	r.Creates++
	return nil
}

// LinkIdentity links the identity, lowercases the account email and, if the
// account had a password, clears it and deletes its reset code.
func (r *IdentityRepository) LinkIdentity(_ context.Context, identity user.LinkedIdentity) error {
	if r.LinkErr != nil {
		return r.LinkErr
	}
	key := identityKey(identity.Provider, identity.Subject)
	if _, dup := r.identities[key]; dup {
		return user.ErrRegistrationFailed
	}

	if account, ok := r.users.accounts[identity.UserID]; ok {
		account.Email = strings.ToLower(account.Email)
		if !user.IsPasswordless(r.users.hashes[account.ID]) {
			r.users.hashes[account.ID] = user.NoPasswordHash
			delete(r.resets.codes, account.ID)
		}
	}
	r.identities[key] = identity
	r.Links++
	return nil
}
