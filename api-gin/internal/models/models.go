// Package models defines the GORM models and DTOs used throughout the application.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model with UUID string ID
type Base struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt time.Time      `                            json:"createdAt"`
	UpdatedAt time.Time      `                            json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

type User struct {
	Base
	Email         string       `gorm:"uniqueIndex;not null"     json:"email"`
	PasswordHash  string       `gorm:"not null"                 json:"-"`
	Name          string       `gorm:"not null"                 json:"name"`
	LastLoginAt   *time.Time   `                                json:"lastLoginAt"`
	UserProfileID *string      `gorm:"type:uuid"                json:"profileId"`
	UserProfile   *UserProfile `gorm:"foreignKey:UserProfileID" json:"profile,omitempty"`
}

type UserProfile struct {
	Base
	Name                     string             `gorm:"not null"                 json:"name"`
	Age                      *int               `                                json:"age"`
	Height                   *float64           `                                json:"height"`
	CurrentWeight            *float64           `                                json:"currentWeight"`
	GoalWeight               *float64           `                                json:"goalWeight"`
	AvatarURL                *string            `                                json:"avatarUrl"`
	ExperienceLevel          string             `gorm:"default:'Beginner'"       json:"experienceLevel"`
	AvailableDaysPerWeek     int                `gorm:"default:3"                json:"availableDaysPerWeek"`
	PreferredWorkoutDuration *int               `                                json:"preferredWorkoutDuration"`
	EquipmentType            string             `gorm:"default:'Gym'"            json:"equipmentType"`
	Workouts                 []Workout          `gorm:"foreignKey:UserProfileID" json:"workouts,omitempty"`
	BioimpedanceHistory      []BioimpedanceData `gorm:"foreignKey:UserProfileID" json:"bioimpedanceHistory,omitempty"`
	WaterIntake              []WaterIntakeEntry `gorm:"foreignKey:UserProfileID" json:"waterIntake,omitempty"`
}

type Workout struct {
	Base
	Name          string     `gorm:"not null"                                         json:"name"`
	Description   *string    `                                                        json:"description"`
	Goal          string     `                                                        json:"goal"`
	Days          string     `                                                        json:"days"` // Comma-separated string for GORM simplicity
	UserProfileID string     `gorm:"type:uuid;not null"                               json:"userProfileId"`
	IsCompleted   bool       `gorm:"default:false"                                    json:"isCompleted"`
	CompletedAt   *time.Time `                                                        json:"completedAt"`
	Exercises     []Exercise `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"exercises"`
}

type Exercise struct {
	Base
	Name        string   `gorm:"not null"           json:"name"`
	MuscleGroup string   `                          json:"muscleGroup"`
	Sets        int      `                          json:"sets"`
	Reps        int      `                          json:"reps"`
	Weight      *float64 `                          json:"weight"`
	RestSeconds int      `                          json:"restSeconds"`
	Notes       *string  `                          json:"notes"`
	WorkoutID   string   `gorm:"type:uuid;not null" json:"workoutId"`
}

type CompletedWorkout struct {
	Base
	WorkoutID       string    `gorm:"type:uuid;not null"        json:"workoutId"`
	UserProfileID   string    `gorm:"type:uuid;not null"        json:"userProfileId"`
	CompletedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"completedAt"`
	DurationSeconds int       `                                 json:"durationSeconds"`
}

type BioimpedanceData struct {
	Base
	Date              time.Time `gorm:"not null"           json:"date"`
	Weight            float64   `                          json:"weight"`
	BodyFatPercentage float64   `                          json:"bodyFatPercentage"`
	MuscleMass        float64   `                          json:"muscleMass"`
	BoneMass          float64   `                          json:"boneMass"`
	WaterPercentage   float64   `                          json:"waterPercentage"`
	VisceralFat       int       `                          json:"visceralFat"`
	BMR               int       `                          json:"bmr"`
	MetabolicAge      int       `                          json:"metabolicAge"`
	Notes             *string   `                          json:"notes"`
	UserProfileID     string    `gorm:"type:uuid;not null" json:"userProfileId"`
}

type WaterIntakeEntry struct {
	Base
	UserProfileID string    `gorm:"type:uuid;not null;index:idx_water_user_date" json:"userProfileId"`
	AmountMl      int       `gorm:"not null"                                     json:"amountMl"`
	ConsumedAt    time.Time `gorm:"not null;index:idx_water_user_date"           json:"consumedAt"`
	Note          *string   `                                                    json:"note"`
}

// DTOs for Auth
type RegisterRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"     binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
