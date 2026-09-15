// Package infrastructure implements the workout slice outbound ports on SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"fittracker-api/internal/workout"
)

// workoutColumns is the column list the previous implementation selected.
const workoutColumns = `id, name, description, goal, days, created_at, updated_at, is_completed, completed_at`

// exerciseColumns is the column list the previous implementation selected.
const exerciseColumns = `id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id`

// WorkoutRepository is the SQLite implementation of workout.Repository.
type WorkoutRepository struct {
	db *sql.DB
}

var _ workout.Repository = (*WorkoutRepository)(nil)

// NewWorkoutRepository builds the workout repository.
func NewWorkoutRepository(db *sql.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

// GetByProfile returns every workout of a profile, newest first, with its
// exercises attached.
func (r *WorkoutRepository) GetByProfile(ctx context.Context, profileID string) ([]workout.Workout, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+workoutColumns+` FROM Workouts WHERE user_profile_id = ? ORDER BY created_at DESC`,
		profileID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts, err := scanWorkouts(rows)
	if err != nil {
		return nil, err
	}

	if err := r.attachExercises(ctx, workouts); err != nil {
		return nil, err
	}

	return workouts, nil
}

// GetByProfileAndDay returns the workouts scheduled for a weekday. Matching is
// comma-delimited so a day name cannot match inside another one.
func (r *WorkoutRepository) GetByProfileAndDay(ctx context.Context, profileID, day string) ([]workout.Workout, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+workoutColumns+` FROM Workouts WHERE user_profile_id = ? AND ',' || LOWER(days) || ',' LIKE '%,' || ? || ',%'`,
		profileID,
		day,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts, err := scanWorkouts(rows)
	if err != nil {
		return nil, err
	}

	if err := r.attachExercises(ctx, workouts); err != nil {
		return nil, err
	}

	return workouts, nil
}

// GetByID returns a workout scoped to its profile, with its exercises attached.
func (r *WorkoutRepository) GetByID(ctx context.Context, id, profileID string) (*workout.Workout, bool, error) {
	var w workout.Workout
	var days string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT `+workoutColumns+` FROM Workouts WHERE id = ? AND user_profile_id = ?`,
		id,
		profileID,
	).Scan(
		&w.ID,
		&w.Name,
		&w.Description,
		&w.Goal,
		&days,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.IsCompleted,
		&w.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	w.Days = strings.Split(days, ",")

	exercises, err := r.loadExercises(ctx, w.ID)
	if err != nil {
		return nil, false, err
	}
	w.Exercises = exercises

	return &w, true, nil
}

// Save creates the workout and its exercises in a single transaction.
func (r *WorkoutRepository) Save(ctx context.Context, profileID string, w *workout.Workout) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO Workouts (id, name, description, goal, days, user_profile_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		w.ID,
		w.Name,
		w.Description,
		w.Goal,
		strings.Join(w.Days, ","),
		profileID,
		now,
		now,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := insertExercises(ctx, tx, w); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Update replaces the workout fields and re-creates its exercises.
func (r *WorkoutRepository) Update(ctx context.Context, id, profileID string, w *workout.Workout) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE Workouts SET name = ?, description = ?, goal = ?, days = ?, updated_at = ? WHERE id = ? AND user_profile_id = ?`,
		w.Name,
		w.Description,
		w.Goal,
		strings.Join(w.Days, ","),
		time.Now(),
		id,
		profileID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Simple approach: drop the old exercises and insert the submitted ones.
	if _, err = tx.ExecContext(ctx, `DELETE FROM Exercises WHERE workout_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := insertExercises(ctx, tx, w); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete removes a workout, cascading to its exercises.
func (r *WorkoutRepository) Delete(ctx context.Context, id, profileID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Workouts WHERE id = ? AND user_profile_id = ?`, id, profileID)
	return err
}

// insertExercises writes every exercise of a workout inside tx.
func insertExercises(ctx context.Context, tx *sql.Tx, w *workout.Workout) error {
	for _, ex := range w.Exercises {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO Exercises (id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			ex.ID,
			ex.Name,
			ex.MuscleGroup,
			ex.Sets,
			ex.Reps,
			ex.Weight,
			ex.RestSeconds,
			ex.Notes,
			w.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// scanWorkouts drains every row into a slice.
func scanWorkouts(rows *sql.Rows) ([]workout.Workout, error) {
	var workouts []workout.Workout

	for rows.Next() {
		var w workout.Workout
		var days string

		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Description,
			&w.Goal,
			&days,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.IsCompleted,
			&w.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		w.Days = strings.Split(days, ",")
		workouts = append(workouts, w)
	}

	return workouts, rows.Err()
}

// loadExercises returns the exercises of a single workout.
func (r *WorkoutRepository) loadExercises(ctx context.Context, workoutID string) ([]workout.Exercise, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+exerciseColumns+` FROM Exercises WHERE workout_id = ?`,
		workoutID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []workout.Exercise
	for rows.Next() {
		var ex workout.Exercise
		err := rows.Scan(
			&ex.ID,
			&ex.Name,
			&ex.MuscleGroup,
			&ex.Sets,
			&ex.Reps,
			&ex.Weight,
			&ex.RestSeconds,
			&ex.Notes,
			&ex.WorkoutID,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}

	return exercises, rows.Err()
}

// attachExercises loads the exercises of every workout in a single query
// instead of one query per workout, and attaches them in place. A workout
// without exercises keeps a nil slice, which the API serialises as JSON null.
func (r *WorkoutRepository) attachExercises(ctx context.Context, workouts []workout.Workout) error {
	if len(workouts) == 0 {
		return nil
	}

	placeholders := make([]string, len(workouts))
	args := make([]any, len(workouts))
	for i := range workouts {
		placeholders[i] = "?"
		args[i] = workouts[i].ID
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+exerciseColumns+` FROM Exercises WHERE workout_id IN (`+strings.Join(placeholders, ",")+`)`,
		args...,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	byWorkout := map[string][]workout.Exercise{}
	for rows.Next() {
		var ex workout.Exercise
		err := rows.Scan(
			&ex.ID,
			&ex.Name,
			&ex.MuscleGroup,
			&ex.Sets,
			&ex.Reps,
			&ex.Weight,
			&ex.RestSeconds,
			&ex.Notes,
			&ex.WorkoutID,
		)
		if err != nil {
			return err
		}
		byWorkout[ex.WorkoutID] = append(byWorkout[ex.WorkoutID], ex)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i := range workouts {
		workouts[i].Exercises = byWorkout[workouts[i].ID]
	}

	return nil
}
