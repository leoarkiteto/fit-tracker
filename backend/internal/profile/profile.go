// Package profile owns the trainee profile domain.
package profile

import (
	"time"

	"fittracker-api/pkg/apperror"
)

// Profile is the training profile a user works out with.
type Profile struct {
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
	AvatarURL                *string   `json:"avatarUrl"`
	Height                   *float64  `json:"height"`
	CurrentWeight            *float64  `json:"currentWeight"`
	GoalWeight               *float64  `json:"goalWeight"`
	PreferredWorkoutDuration *int      `json:"preferredWorkoutDuration"`
	Age                      *int      `json:"age"`
	ID                       string    `json:"id"`
	ExperienceLevel          string    `json:"experienceLevel"`
	EquipmentType            string    `json:"equipmentType"`
	Name                     string    `json:"name"`
	AvailableDaysPerWeek     int       `json:"availableDaysPerWeek"`
}

// ErrProfileNotFound is returned when a profile id does not exist.
var ErrProfileNotFound = apperror.NotFound("Profile not found")
