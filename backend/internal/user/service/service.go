// Package service holds the account use cases. It depends only on the domain
// and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/user"
	"fittracker-api/pkg/auth"
)

// Service holds the account use cases.
type Service struct {
	repo   user.Repository
	signer *auth.Signer
}

// New builds the account service.
func New(repo user.Repository, signer *auth.Signer) *Service {
	return &Service{repo: repo, signer: signer}
}
