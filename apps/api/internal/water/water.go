// Package water owns the hydration-tracking domain.
package water

import (
	"time"

	"fittracker-api/pkg/apperror"
)

// Entry is one logged drink.
type Entry struct {
	ConsumedAt    time.Time `json:"consumedAt"`
	Note          *string   `json:"note"`
	ID            string    `json:"id"`
	UserProfileID string    `json:"userProfileId"`
	AmountMl      int       `json:"amountMl"`
}

// DailySummary is the hydration picture for a single day.
type DailySummary struct {
	Date    string  `json:"date"`
	TotalMl int     `json:"totalMl"`
	GoalMl  int     `json:"goalMl"`
	Entries []Entry `json:"entries"`
}

// The daily target is 35 ml per kg of body weight, using 70 kg when no weight
// is recorded for the profile.
const (
	MlPerKg         = 35
	DefaultWeightKg = 70.0
)

// ErrInvalidDate is returned when the ?date= parameter cannot be parsed.
var ErrInvalidDate = apperror.BadRequest("Invalid date format")
