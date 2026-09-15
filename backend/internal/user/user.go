// Package user owns the account and authentication domain.
package user

import (
	"net/http"
	"time"

	"fittracker-api/pkg/apperror"
)

// User is an account that can authenticate against the API.
type User struct {
	CreatedAt    time.Time  `json:"createdAt"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
	ProfileID    *string    `json:"profileId"`
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         string     `json:"name"`
}

// DTO is the public projection of a User.
type DTO struct {
	ProfileID *string `json:"profileId"`
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
}

// ToDTO projects the account onto its public shape.
func (u *User) ToDTO() DTO {
	return DTO{ProfileID: u.ProfileID, ID: u.ID, Email: u.Email, Name: u.Name}
}

// AuthResponse is the session returned by register and login.
type AuthResponse struct {
	User      DTO       `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token"`
}

// Domain errors. The messages are part of the API contract the mobile client
// displays, so they are defined next to the domain rather than in the handler.
var (
	// ErrInvalidCredentials covers both an unknown email and a bad password.
	ErrInvalidCredentials = apperror.Unauthorized("Invalid credentials")
	// ErrRegistrationFailed is deliberately generic to avoid email enumeration.
	ErrRegistrationFailed = apperror.Conflict("Registration failed")
	// ErrUserNotFound is returned when an account id does not exist.
	ErrUserNotFound = apperror.NotFound("User not found")
	// ErrNewPasswordRequired is returned for an empty new password.
	ErrNewPasswordRequired = apperror.BadRequest("New password is required")
	// ErrCurrentPasswordMismatch is returned when the current password is wrong.
	ErrCurrentPasswordMismatch = apperror.Unauthorized("Current password is incorrect")
	// ErrHashingPassword wraps a bcrypt failure during registration.
	ErrHashingPassword = apperror.New(http.StatusInternalServerError, "Error hashing password")
	// ErrGeneratingToken wraps a signing failure.
	ErrGeneratingToken = apperror.New(http.StatusInternalServerError, "Error generating token")
)
