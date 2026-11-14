// Package bioimpedance owns the body-composition measurement domain.
package bioimpedance

import (
	"time"

	"fittracker-api/pkg/apperror"
)

// Measurement is a single bioimpedance reading.
type Measurement struct {
	Date              time.Time `json:"date"`
	Notes             *string   `json:"notes"`
	ID                string    `json:"id"`
	UserProfileID     string    `json:"userProfileId"`
	Weight            float64   `json:"weight"`
	BodyFatPercentage float64   `json:"bodyFatPercentage"`
	MuscleMass        float64   `json:"muscleMass"`
	BoneMass          float64   `json:"boneMass"`
	WaterPercentage   float64   `json:"waterPercentage"`
	VisceralFat       int       `json:"visceralFat"`
	BMR               int       `json:"bmr"`
	MetabolicAge      int       `json:"metabolicAge"`
}

// ErrNoMeasurement is returned when a profile has no reading yet.
var ErrNoMeasurement = apperror.NotFound("No data found")
