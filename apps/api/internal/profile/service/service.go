// Package service holds the profile use cases. It depends only on the domain
// and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/profile"
)

// Service holds the profile use cases.
type Service struct {
	repo profile.Repository
}

// New builds the profile service.
func New(repo profile.Repository) *Service {
	return &Service{repo: repo}
}
