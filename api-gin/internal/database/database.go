// Package database handles the GORM database connection, initialization, and auto-migrations.
package database

import (
	"log"

	"fittracker-api-gin/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dataSourceName string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Workout{},
		&models.Exercise{},
		&models.CompletedWorkout{},
		&models.BioimpedanceData{},
		&models.WaterIntakeEntry{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return DB
}
