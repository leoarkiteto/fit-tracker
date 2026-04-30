package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"fittracker-api/internal/models"
	"github.com/google/uuid"
)

type WorkoutHandler struct {
	DB *sql.DB
}

func (h *WorkoutHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")

	rows, err := h.DB.Query(`SELECT id, name, description, goal, days, created_at, updated_at, is_completed, completed_at FROM Workouts WHERE user_profile_id = ? ORDER BY created_at DESC`, profileID)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var workouts []models.Workout
	var workoutIDs []string
	for rows.Next() {
		var workout models.Workout
		var daysStr string
		err := rows.Scan(&workout.ID, &workout.Name, &workout.Description, &workout.Goal, &daysStr, &workout.CreatedAt, &workout.UpdatedAt, &workout.IsCompleted, &workout.CompletedAt)
		if err != nil {
			serverError(w, err)
			return
		}
		workout.Days = strings.Split(daysStr, ",")
		workout.UserProfileID = profileID
		workoutIDs = append(workoutIDs, workout.ID)
		workouts = append(workouts, workout)
	}
	rows.Close()

	// Batch load exercises (N+1 → 1 query)
	exerciseMap, err := h.loadExercisesBatch(workoutIDs)
	if err != nil {
		serverError(w, err)
		return
	}
	for i := range workouts {
		workouts[i].Exercises = exerciseMap[workouts[i].ID]
	}

	JSONResponse(w, http.StatusOK, workouts)
}

func (h *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profileID := r.PathValue("profileId")

	var workout models.Workout
	var daysStr string
	err := h.DB.QueryRow(`SELECT id, name, description, goal, days, created_at, updated_at, is_completed, completed_at FROM Workouts WHERE id = ? AND user_profile_id = ?`, id, profileID).
		Scan(&workout.ID, &workout.Name, &workout.Description, &workout.Goal, &daysStr, &workout.CreatedAt, &workout.UpdatedAt, &workout.IsCompleted, &workout.CompletedAt)

	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	workout.Days = strings.Split(daysStr, ",")
	workout.UserProfileID = profileID
	
	exercises, err := h.loadExercises(workout.ID)
	if err != nil {
		serverError(w, err)
		return
	}
	workout.Exercises = exercises

	JSONResponse(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	var req models.Workout
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workoutID := uuid.New().String()
	daysStr := strings.Join(req.Days, ",")

	tx, err := h.DB.Begin()
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = tx.Exec(`INSERT INTO Workouts (id, name, description, goal, days, user_profile_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		workoutID, req.Name, req.Description, req.Goal, daysStr, profileID, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}

	for _, ex := range req.Exercises {
		exID := uuid.New().String()
		_, err = tx.Exec(`INSERT INTO Exercises (id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			exID, ex.Name, ex.MuscleGroup, ex.Sets, ex.Reps, ex.Weight, ex.RestSeconds, ex.Notes, workoutID)
		if err != nil {
			tx.Rollback()
			serverError(w, err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	req.ID = workoutID
	JSONResponse(w, http.StatusCreated, req)
}

func (h *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profileID := r.PathValue("profileId")
	var req models.Workout
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	daysStr := strings.Join(req.Days, ",")

	tx, err := h.DB.Begin()
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = tx.Exec(`UPDATE Workouts SET name = ?, description = ?, goal = ?, days = ?, updated_at = ? WHERE id = ? AND user_profile_id = ?`,
		req.Name, req.Description, req.Goal, daysStr, time.Now(), id, profileID)
	if err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}

	// Simple way: delete old exercises and insert new ones
	_, err = tx.Exec(`DELETE FROM Exercises WHERE workout_id = ?`, id)
	if err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}

	for _, ex := range req.Exercises {
		exID := uuid.New().String()
		_, err = tx.Exec(`INSERT INTO Exercises (id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			exID, ex.Name, ex.MuscleGroup, ex.Sets, ex.Reps, ex.Weight, ex.RestSeconds, ex.Notes, id)
		if err != nil {
			tx.Rollback()
			serverError(w, err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	req.ID = id
	JSONResponse(w, http.StatusOK, req)
}

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profileID := r.PathValue("profileId")

	_, err := h.DB.Exec(`DELETE FROM Workouts WHERE id = ? AND user_profile_id = ?`, id, profileID)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	// Use lowercase English day for matching against comma-separated days column
	today := strings.ToLower(time.Now().Weekday().String())

	// Use comma-prefixed/suffixed matching to avoid substring false positives
	rows, err := h.DB.Query(`SELECT id, name, description, goal, days, created_at, updated_at, is_completed, completed_at FROM Workouts WHERE user_profile_id = ? AND ',' || LOWER(days) || ',' LIKE '%,' || ? || ',%'`, profileID, today)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var workouts []models.Workout
	var workoutIDs []string
	for rows.Next() {
		var workout models.Workout
		var daysStr string
		err := rows.Scan(&workout.ID, &workout.Name, &workout.Description, &workout.Goal, &daysStr, &workout.CreatedAt, &workout.UpdatedAt, &workout.IsCompleted, &workout.CompletedAt)
		if err != nil {
			serverError(w, err)
			return
		}
		workout.Days = strings.Split(daysStr, ",")
		workoutIDs = append(workoutIDs, workout.ID)
		workouts = append(workouts, workout)
	}
	rows.Close()

	// Batch load exercises (N+1 → 1 query)
	exerciseMap, err := h.loadExercisesBatch(workoutIDs)
	if err != nil {
		serverError(w, err)
		return
	}
	for i := range workouts {
		workouts[i].Exercises = exerciseMap[workouts[i].ID]
	}

	JSONResponse(w, http.StatusOK, workouts)
}

func (h *WorkoutHandler) loadExercises(workoutID string) ([]models.Exercise, error) {
	rows, err := h.DB.Query(`SELECT id, name, muscle_group, sets, reps, weight, rest_seconds, notes FROM Exercises WHERE workout_id = ?`, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		err := rows.Scan(&ex.ID, &ex.Name, &ex.MuscleGroup, &ex.Sets, &ex.Reps, &ex.Weight, &ex.RestSeconds, &ex.Notes)
		if err != nil {
			return nil, err
		}
		ex.WorkoutID = workoutID
		exercises = append(exercises, ex)
	}
	return exercises, nil
}

// loadExercisesBatch loads exercises for multiple workout IDs in a single query.
// Returns a map of workoutID → []Exercise to eliminate the N+1 problem.
func (h *WorkoutHandler) loadExercisesBatch(workoutIDs []string) (map[string][]models.Exercise, error) {
	result := make(map[string][]models.Exercise)
	if len(workoutIDs) == 0 {
		return result, nil
	}

	// Build placeholders (?, ?, ...)
	placeholders := make([]string, len(workoutIDs))
	args := make([]interface{}, len(workoutIDs))
	for i, id := range workoutIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id FROM Exercises WHERE workout_id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ex models.Exercise
		err := rows.Scan(&ex.ID, &ex.Name, &ex.MuscleGroup, &ex.Sets, &ex.Reps, &ex.Weight, &ex.RestSeconds, &ex.Notes, &ex.WorkoutID)
		if err != nil {
			return nil, err
		}
		result[ex.WorkoutID] = append(result[ex.WorkoutID], ex)
	}

	return result, rows.Err()
}
