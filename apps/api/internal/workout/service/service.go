// Package service holds the workout use cases. It depends only on the domain
// and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/workout"
)

// Service holds the workout use cases.
type Service struct {
	repo workout.Repository
}

// New builds the workout service.
func New(repo workout.Repository) *Service {
	return &Service{repo: repo}
}
