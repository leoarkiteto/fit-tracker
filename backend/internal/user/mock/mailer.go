package mock

import (
	"context"

	"fittracker-api/internal/user"
)

// SentCode records one delivered reset code.
type SentCode struct {
	Email string
	Code  string
}

// Mailer is a recording user.Mailer. Set SendErr to force delivery to fail.
type Mailer struct {
	SendErr error

	Sent []SentCode
}

var _ user.Mailer = (*Mailer)(nil)

// NewMailer builds a mailer with nothing sent.
func NewMailer() *Mailer { return &Mailer{} }

// SendPasswordResetCode records the delivery unless SendErr is set.
func (m *Mailer) SendPasswordResetCode(_ context.Context, email, code string) error {
	if m.SendErr != nil {
		return m.SendErr
	}
	m.Sent = append(m.Sent, SentCode{Email: email, Code: code})
	return nil
}
