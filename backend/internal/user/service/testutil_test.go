package service

import (
	"testing"

	"fittracker-api/internal/user"
	"fittracker-api/pkg/auth"
)

// newTestService builds the service under test with an in-memory repository and
// a signer that does not depend on the JWT_SECRET environment variable.
func newTestService(t *testing.T, repo user.Repository) (*Service, *auth.Signer) {
	t.Helper()

	signer, err := auth.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("building signer: %v", err)
	}

	return New(repo, signer), signer
}
