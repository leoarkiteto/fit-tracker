package models

import (
	"time"
)

// Auth
type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Name          string     `json:"name"`
	CreatedAt     time.Time  `json:"createdAt"`
	LastLoginAt   *time.Time `json:"lastLoginAt"`
	UserProfileID *string    `json:"profileId"`
}

// Profiles
type UserProfile struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	Age                      *int    `json:"age"`
	Height                   *float64 `json:"height"`
	CurrentWeight            *float64 `json:"currentWeight"`
	GoalWeight               *float64 `json:"goalWeight"`
	AvatarURL                *string `json:"avatarUrl"`
	ExperienceLevel          string  `json:"experienceLevel"`
	AvailableDaysPerWeek     int     `json:"availableDaysPerWeek"`
	PreferredWorkoutDuration *int    `json:"preferredWorkoutDuration"`
	EquipmentType            string  `json:"equipmentType"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

// Workouts
type Workout struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	Goal          string     `json:"goal"`
	Days          []string   `json:"days"`
	UserProfileID string     `json:"userProfileId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	IsCompleted   bool       `json:"isCompleted"`
	CompletedAt   *time.Time `json:"completedAt"`
	Exercises     []Exercise `json:"exercises"`
}

type Exercise struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MuscleGroup string   `json:"muscleGroup"`
	Sets        int      `json:"sets"`
	Reps        int      `json:"reps"`
	Weight      *float64 `json:"weight"`
	RestSeconds int      `json:"restSeconds"`
	Notes       *string  `json:"notes"`
	WorkoutID   string   `json:"workoutId"`
}

type CompletedWorkout struct {
	ID              string    `json:"id"`
	WorkoutID       string    `json:"workoutId"`
	UserProfileID   string    `json:"userProfileId"`
	CompletedAt     time.Time `json:"completedAt"`
	DurationSeconds int       `json:"durationSeconds"`
}

// Bioimpedance
type BioimpedanceData struct {
	ID                string    `json:"id"`
	Date              time.Time `json:"date"`
	Weight            float64   `json:"weight"`
	BodyFatPercentage float64   `json:"bodyFatPercentage"`
	MuscleMass        float64   `json:"muscleMass"`
	BoneMass          float64   `json:"boneMass"`
	WaterPercentage   float64   `json:"waterPercentage"`
	VisceralFat       int       `json:"visceralFat"`
	BMR               int       `json:"bmr"`
	MetabolicAge      int       `json:"metabolicAge"`
	Notes             *string   `json:"notes"`
	UserProfileID     string    `json:"userProfileId"`
}

// Water Intake
type WaterIntakeEntry struct {
	ID            string    `json:"id"`
	UserProfileID string    `json:"userProfileId"`
	AmountMl      int       `json:"amountMl"`
	ConsumedAt    time.Time `json:"consumedAt"`
	Note          *string   `json:"note"`
}

// DTOs for Auth
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
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      UserDto   `json:"user"`
}

type UserDto struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	ProfileID *string `json:"profileId"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
