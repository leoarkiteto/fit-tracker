package infrastructure

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/aiplanning"
)

// WorkoutWriter stores accepted plans in the Workouts and Exercises tables.
//
// It intentionally keeps its own copy of the two INSERT statements rather than
// depending on the workout slice, so the AI feature stays a self-contained
// vertical slice.
type WorkoutWriter struct {
	db *sql.DB
}

var _ aiplanning.WorkoutWriter = (*WorkoutWriter)(nil)

// NewWorkoutWriter builds the plan writer.
func NewWorkoutWriter(db *sql.DB) *WorkoutWriter {
	return &WorkoutWriter{db: db}
}

// CreateAll stores every workout of a plan for a profile in a single
// transaction, so an accepted plan is all-or-nothing.
func (w *WorkoutWriter) CreateAll(
	ctx context.Context,
	profileID string,
	planned []aiplanning.PlannedWorkout,
) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, workout := range planned {
		if err := insertPlannedWorkout(ctx, tx, profileID, workout); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// insertPlannedWorkout writes one workout and its exercises inside tx.
func insertPlannedWorkout(
	ctx context.Context,
	tx *sql.Tx,
	profileID string,
	planned aiplanning.PlannedWorkout,
) error {
	workoutID := uuid.New().String()
	now := time.Now()

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO Workouts (id, name, description, goal, days, user_profile_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		workoutID,
		planned.Name,
		planned.Description,
		planned.Goal,
		strings.Join(planned.Days, ","),
		profileID,
		now,
		now,
	)
	if err != nil {
		return err
	}

	for _, exercise := range planned.Exercises {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO Exercises (id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			uuid.New().String(),
			exercise.Name,
			exercise.MuscleGroup,
			exercise.Sets,
			exercise.Reps,
			exercise.Weight,
			exercise.RestSeconds,
			exercise.Notes,
			workoutID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
