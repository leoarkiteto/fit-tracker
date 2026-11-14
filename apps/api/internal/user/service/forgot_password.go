package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/apperror"
)

// hashResetCode returns the hex HMAC-SHA256 of code under the reset key. Only
// this value is ever persisted.
func (s *Service) hashResetCode(code string) string {
	mac := hmac.New(sha256.New, s.resetKey)
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

// RequestPasswordReset issues a reset code for the account with the given email
// and delivers it. It returns nil whether or not the account exists, and inside
// the resend cooldown, so callers cannot enumerate accounts.
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	account, found, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("request password reset: loading account: %v", err)
		return apperror.Internal()
	}
	if !found {
		return nil
	}

	now := s.now().UTC()

	existing, exists, err := s.resets.Get(ctx, account.ID)
	if err != nil {
		log.Printf("request password reset: loading code: %v", err)
		return apperror.Internal()
	}
	if exists && now.Sub(existing.CreatedAt) < user.ResetCodeCooldown {
		return nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		log.Printf("request password reset: generating code: %v", err)
		return apperror.Internal()
	}
	code := fmt.Sprintf("%06d", n.Int64())

	err = s.resets.Replace(ctx, user.PasswordResetCode{
		UserID:    account.ID,
		CodeHash:  s.hashResetCode(code),
		CreatedAt: now,
		Attempts:  0,
	})
	if err != nil {
		log.Printf("request password reset: storing code: %v", err)
		return apperror.Internal()
	}

	if err := s.mailer.SendPasswordResetCode(ctx, account.Email, code); err != nil {
		log.Printf("request password reset: delivering code: %v", err)
	}

	return nil
}
