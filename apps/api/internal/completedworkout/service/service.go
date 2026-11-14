// Package service holds the completed-workout use cases. It depends only on the
// domain and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/completedworkout"
)

// Service holds the completed-workout use cases.
type Service struct {
	repo completedworkout.Repository
}

// New builds the completed-workout service.
func New(repo completedworkout.Repository) *Service {
	return &Service{repo: repo}
}
