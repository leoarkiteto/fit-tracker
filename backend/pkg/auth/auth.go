// Package auth provides password hashing and JWT issuing/validation primitives
// shared by every slice that needs to authenticate a caller.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ErrMissingSecret is returned when no JWT secret is configured.
var ErrMissingSecret = errors.New("JWT_SECRET environment variable is not set")

const tokenTTL = 24 * time.Hour

// Claims are the values embedded into an issued token.
type Claims struct {
	UserID    string
	Email     string
	Name      string
	ProfileID string
}

// Signer issues and validates JWTs with a single secret.
type Signer struct {
	secret []byte
}

// NewSigner builds a Signer, refusing to start without a secret.
func NewSigner(secret string) (*Signer, error) {
	if secret == "" {
		return nil, ErrMissingSecret
	}
	return &Signer{secret: []byte(secret)}, nil
}

// GenerateToken signs a token for the given claims and reports when it expires.
func (s *Signer) GenerateToken(c Claims) (string, time.Time, error) {
	expiresAt := time.Now().Add(tokenTTL)

	claims := jwt.MapClaims{
		"sub":       c.UserID,
		"email":     c.Email,
		"name":      c.Name,
		"profileId": c.ProfileID,
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken parses a token and returns its claims when it is valid.
func (s *Signer) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash reports whether password matches the stored bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
