// Package models
package models

import (
	"time"
)

// User type
type User struct {
	CreatedAt     time.Time  `json:"createdAt"`
	LastLoginAt   *time.Time `json:"lastLoginAt"`
	UserProfileID *string    `json:"profileId"`
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Name          string     `json:"name"`
}

// UserProfile type
type UserProfile struct {
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

// Workout type
type Workout struct {
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Description   *string    `json:"description"`
	CompletedAt   *time.Time `json:"completedAt"`
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Goal          string     `json:"goal"`
	UserProfileID string     `json:"userProfileId"`
	Days          []string   `json:"days"`
	Exercises     []Exercise `json:"exercises"`
	IsCompleted   bool       `json:"isCompleted"`
}

type Exercise struct {
	Weight      *float64 `json:"weight"`
	Notes       *string  `json:"notes"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MuscleGroup string   `json:"muscleGroup"`
	WorkoutID   string   `json:"workoutId"`
	Sets        int      `json:"sets"`
	Reps        int      `json:"reps"`
	RestSeconds int      `json:"restSeconds"`
}

type CompletedWorkout struct {
	CompletedAt     time.Time `json:"completedAt"`
	ID              string    `json:"id"`
	WorkoutID       string    `json:"workoutId"`
	UserProfileID   string    `json:"userProfileId"`
	DurationSeconds int       `json:"durationSeconds"`
}

// BioimpedanceData type
type BioimpedanceData struct {
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

// WaterIntakeEntry type
type WaterIntakeEntry struct {
	ConsumedAt    time.Time `json:"consumedAt"`
	Note          *string   `json:"note"`
	ID            string    `json:"id"`
	UserProfileID string    `json:"userProfileId"`
	AmountMl      int       `json:"amountMl"`
}

// RegisterRequest type
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User      UserDto   `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token"`
}

type UserDto struct {
	ProfileID *string `json:"profileId"`
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
