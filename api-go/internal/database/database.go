package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"log"
	"time"
)

var DB *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Enable foreign key enforcement (disabled by default in SQLite)
	if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Error enabling foreign keys: %v", err)
	}

	// Connection pool configuration
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	createSchema()
	return DB
}

func createSchema() {
	schema := `
	CREATE TABLE IF NOT EXISTS UserProfiles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INTEGER,
		height REAL,
		current_weight REAL,
		goal_weight REAL,
		avatar_url TEXT,
		experience_level TEXT DEFAULT 'Beginner',
		available_days_per_week INTEGER DEFAULT 3,
		preferred_workout_duration INTEGER,
		equipment_type TEXT DEFAULT 'Gym',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS Users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME,
		user_profile_id TEXT,
		FOREIGN KEY (user_profile_id) REFERENCES UserProfiles(id)
	);

	CREATE TABLE IF NOT EXISTS Workouts (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		goal TEXT,
		days TEXT, -- Stored as comma-separated values
		user_profile_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_completed BOOLEAN DEFAULT 0,
		completed_at DATETIME,
		FOREIGN KEY (user_profile_id) REFERENCES UserProfiles(id)
	);

	CREATE TABLE IF NOT EXISTS Exercises (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		muscle_group TEXT,
		sets INTEGER,
		reps INTEGER,
		weight REAL,
		rest_seconds INTEGER,
		notes TEXT,
		workout_id TEXT NOT NULL,
		FOREIGN KEY (workout_id) REFERENCES Workouts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS CompletedWorkouts (
		id TEXT PRIMARY KEY,
		workout_id TEXT NOT NULL,
		user_profile_id TEXT NOT NULL,
		completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		duration_seconds INTEGER,
		FOREIGN KEY (workout_id) REFERENCES Workouts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_profile_id) REFERENCES UserProfiles(id)
	);

	CREATE TABLE IF NOT EXISTS BioimpedanceData (
		id TEXT PRIMARY KEY,
		date DATETIME NOT NULL,
		weight REAL,
		body_fat_percentage REAL,
		muscle_mass REAL,
		bone_mass REAL,
		water_percentage REAL,
		visceral_fat INTEGER,
		bmr INTEGER,
		metabolic_age INTEGER,
		notes TEXT,
		user_profile_id TEXT NOT NULL,
		FOREIGN KEY (user_profile_id) REFERENCES UserProfiles(id)
	);

	CREATE TABLE IF NOT EXISTS WaterIntakeEntries (
		id TEXT PRIMARY KEY,
		user_profile_id TEXT NOT NULL,
		amount_ml INTEGER NOT NULL,
		consumed_at DATETIME NOT NULL,
		note TEXT,
		FOREIGN KEY (user_profile_id) REFERENCES UserProfiles(id)
	);
	
	CREATE INDEX IF NOT EXISTS IX_WaterIntakeEntries_UserProfileId_ConsumedAt ON WaterIntakeEntries (user_profile_id, consumed_at);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}
}
