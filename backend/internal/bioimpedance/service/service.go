// Package service holds the bioimpedance use cases. It depends only on the
// domain and its ports, so it can be unit tested without HTTP or a database.
package service

import (
	"fittracker-api/internal/bioimpedance"
)

// Service holds the bioimpedance use cases.
type Service struct {
	repo bioimpedance.Repository
}

// New builds the bioimpedance service.
func New(repo bioimpedance.Repository) *Service {
	return &Service{repo: repo}
}
