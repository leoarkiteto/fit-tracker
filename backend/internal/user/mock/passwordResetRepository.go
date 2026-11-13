package mock

import (
	"context"
	"time"

	"fittracker-api/internal/user"
)

// PasswordResetRepository is an in-memory user.PasswordResetRepository. Set any
// of the Err fields to force that call to fail. ConsumeAttempt and MarkUsed
// mirror the conditional SQL updates of the real adapter.
type PasswordResetRepository struct {
	ReplaceErr        error
	GetErr            error
	ConsumeAttemptErr error
	MarkUsedErr       error

	Replaces int

	codes map[string]user.PasswordResetCode
}

var _ user.PasswordResetRepository = (*PasswordResetRepository)(nil)

// NewPasswordResetRepository builds an empty in-memory repository.
func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{codes: map[string]user.PasswordResetCode{}}
}

// Seed stores a code directly, bypassing Replace and its error/counter.
func (r *PasswordResetRepository) Seed(code user.PasswordResetCode) {
	r.codes[code.UserID] = code
}

// Code returns the stored code of an account.
func (r *PasswordResetRepository) Code(userID string) (user.PasswordResetCode, bool) {
	code, ok := r.codes[userID]
	return code, ok
}

// Replace overwrites the account's code.
func (r *PasswordResetRepository) Replace(_ context.Context, code user.PasswordResetCode) error {
	if r.ReplaceErr != nil {
		return r.ReplaceErr
	}
	r.codes[code.UserID] = code
	r.Replaces++
	return nil
}

// Get returns a copy of the account's code.
func (r *PasswordResetRepository) Get(_ context.Context, userID string) (*user.PasswordResetCode, bool, error) {
	if r.GetErr != nil {
		return nil, false, r.GetErr
	}
	code, ok := r.codes[userID]
	if !ok {
		return nil, false, nil
	}
	return &code, true, nil
}

// ConsumeAttempt counts one attempt when the stored code matches codeHash, is
// unused and has fewer than max attempts.
func (r *PasswordResetRepository) ConsumeAttempt(_ context.Context, userID, codeHash string, max int) (bool, error) {
	if r.ConsumeAttemptErr != nil {
		return false, r.ConsumeAttemptErr
	}
	code, ok := r.codes[userID]
	if !ok || code.CodeHash != codeHash || code.UsedAt != nil || code.Attempts >= max {
		return false, nil
	}
	code.Attempts++
	r.codes[userID] = code
	return true, nil
}

// MarkUsed stamps the code as used when it matches codeHash and is still unused.
func (r *PasswordResetRepository) MarkUsed(_ context.Context, userID, codeHash string, at time.Time) (bool, error) {
	if r.MarkUsedErr != nil {
		return false, r.MarkUsedErr
	}
	code, ok := r.codes[userID]
	if !ok || code.CodeHash != codeHash || code.UsedAt != nil {
		return false, nil
	}
	code.UsedAt = &at
	r.codes[userID] = code
	return true, nil
}
