// Package infrastructure implements the completed-workout slice outbound ports
// on SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"fittracker-api/internal/completedworkout"
)

// CompletedWorkoutRepository is the SQLite implementation of
// completedworkout.Repository.
type CompletedWorkoutRepository struct {
	db *sql.DB
}

var _ completedworkout.Repository = (*CompletedWorkoutRepository)(nil)

// NewCompletedWorkoutRepository builds the training-history repository.
func NewCompletedWorkoutRepository(db *sql.DB) *CompletedWorkoutRepository {
	return &CompletedWorkoutRepository{db: db}
}

// GetByProfile returns the history of a profile, newest first.
func (r *CompletedWorkoutRepository) GetByProfile(
	ctx context.Context,
	profileID string,
) ([]completedworkout.CompletedWorkout, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, workout_id, completed_at, duration_seconds FROM CompletedWorkouts WHERE user_profile_id = ? ORDER BY completed_at DESC`,
		profileID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []completedworkout.CompletedWorkout
	for rows.Next() {
		var c completedworkout.CompletedWorkout
		err := rows.Scan(&c.ID, &c.WorkoutID, &c.CompletedAt, &c.DurationSeconds)
		if err != nil {
			return nil, err
		}
		history = append(history, c)
	}

	return history, rows.Err()
}

// Save records a completed workout.
func (r *CompletedWorkoutRepository) Save(ctx context.Context, c *completedworkout.CompletedWorkout) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO CompletedWorkouts (id, workout_id, user_profile_id, completed_at, duration_seconds) VALUES (?, ?, ?, ?, ?)`,
		c.ID,
		c.WorkoutID,
		c.UserProfileID,
		c.CompletedAt,
		c.DurationSeconds,
	)
	return err
}

// CountByProfile counts every completed workout of a profile.
func (r *CompletedWorkoutRepository) CountByProfile(ctx context.Context, profileID string) (int, error) {
	var total int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM CompletedWorkouts WHERE user_profile_id = ?`,
		profileID,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// CountByProfileSince counts the workouts completed at or after an instant.
func (r *CompletedWorkoutRepository) CountByProfileSince(
	ctx context.Context,
	profileID string,
	since time.Time,
) (int, error) {
	var total int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM CompletedWorkouts WHERE user_profile_id = ? AND completed_at >= ?`,
		profileID,
		since,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// SumDurationByProfile totals the recorded duration in seconds.
func (r *CompletedWorkoutRepository) SumDurationByProfile(ctx context.Context, profileID string) (int, error) {
	var seconds int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(duration_seconds), 0) FROM CompletedWorkouts WHERE user_profile_id = ?`,
		profileID,
	).Scan(&seconds)
	if err != nil {
		return 0, err
	}

	return seconds, nil
}
