// Package apperror carries a domain failure together with the HTTP status it
// maps to, so services never have to touch net/http.
package apperror

import (
	"errors"
	"log"
	"net/http"
)

// Error is a domain error with the HTTP status it should be rendered as.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

// New builds an Error with an explicit status.
func New(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

// BadRequest maps to 400.
func BadRequest(message string) *Error { return New(http.StatusBadRequest, message) }

// Unauthorized maps to 401.
func Unauthorized(message string) *Error { return New(http.StatusUnauthorized, message) }

// NotFound maps to 404.
func NotFound(message string) *Error { return New(http.StatusNotFound, message) }

// Conflict maps to 409.
func Conflict(message string) *Error { return New(http.StatusConflict, message) }

// ServiceUnavailable maps to 503.
func ServiceUnavailable(message string) *Error { return New(http.StatusServiceUnavailable, message) }

// Internal is the generic 500 that never leaks implementation details.
func Internal() *Error { return New(http.StatusInternalServerError, "Internal server error") }

// Write renders err as an HTTP error response, falling back to a generic 500
// for anything that is not an *Error.
func Write(w http.ResponseWriter, err error) {
	var appErr *Error
	if errors.As(err, &appErr) {
		if appErr.Status >= http.StatusInternalServerError {
			log.Printf("Internal server error: %v", appErr)
		}
		http.Error(w, appErr.Message, appErr.Status)
		return
	}

	log.Printf("Internal server error: %v", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
