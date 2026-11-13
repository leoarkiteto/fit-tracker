package infrastructure

import (
	"context"
	"log"

	"fittracker-api/internal/user"
)

// LogMailer is a development user.Mailer that writes the email address and
// code to the server log instead of sending anything.
type LogMailer struct{}

var _ user.Mailer = (*LogMailer)(nil)

// NewLogMailer builds the development mailer.
func NewLogMailer() *LogMailer {
	return &LogMailer{}
}

// SendPasswordResetCode logs the reset code for email.
func (m *LogMailer) SendPasswordResetCode(_ context.Context, email, code string) error {
	log.Printf("[DEV MAILER] password reset code for %s: %s", email, code)
	return nil
}
