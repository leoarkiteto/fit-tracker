// Package service holds the hydration use cases. It depends only on the domain
// and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/water"
)

// Service holds the hydration use cases.
type Service struct {
	repo water.Repository
}

// New builds the hydration service.
func New(repo water.Repository) *Service {
	return &Service{repo: repo}
}
