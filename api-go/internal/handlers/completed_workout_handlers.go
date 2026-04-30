package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"fittracker-api/internal/models"
	"github.com/google/uuid"
)

type CompletedWorkoutHandler struct {
	DB *sql.DB
}

func (h *CompletedWorkoutHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")

	rows, err := h.DB.Query(`SELECT id, workout_id, completed_at, duration_seconds FROM CompletedWorkouts WHERE user_profile_id = ? ORDER BY completed_at DESC`, profileID)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var results []models.CompletedWorkout
	for rows.Next() {
		var c models.CompletedWorkout
		err := rows.Scan(&c.ID, &c.WorkoutID, &c.CompletedAt, &c.DurationSeconds)
		if err != nil {
			serverError(w, err)
			return
		}
		c.UserProfileID = profileID
		results = append(results, c)
	}

	JSONResponse(w, http.StatusOK, results)
}

func (h *CompletedWorkoutHandler) Complete(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	var req struct {
		WorkoutID       string `json:"workoutId"`
		DurationSeconds int    `json:"durationSeconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	now := time.Now()

	_, err := h.DB.Exec(`INSERT INTO CompletedWorkouts (id, workout_id, user_profile_id, completed_at, duration_seconds) VALUES (?, ?, ?, ?, ?)`,
		id, req.WorkoutID, profileID, now, req.DurationSeconds)
	if err != nil {
		serverError(w, err)
		return
	}

	JSONResponse(w, http.StatusCreated, models.CompletedWorkout{
		ID:              id,
		WorkoutID:       req.WorkoutID,
		UserProfileID:   profileID,
		CompletedAt:     now,
		DurationSeconds: req.DurationSeconds,
	})
}

func (h *CompletedWorkoutHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")

	var totalWorkouts int
	err := h.DB.QueryRow(`SELECT COUNT(*) FROM CompletedWorkouts WHERE user_profile_id = ?`, profileID).Scan(&totalWorkouts)
	if err != nil {
		serverError(w, err)
		return
	}

	var workoutsThisWeek int
	lastWeek := time.Now().AddDate(0, 0, -7)
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM CompletedWorkouts WHERE user_profile_id = ? AND completed_at >= ?`, profileID, lastWeek).Scan(&workoutsThisWeek)
	if err != nil {
		serverError(w, err)
		return
	}

	var totalSeconds int
	err = h.DB.QueryRow(`SELECT COALESCE(SUM(duration_seconds), 0) FROM CompletedWorkouts WHERE user_profile_id = ?`, profileID).Scan(&totalSeconds)
	if err != nil {
		serverError(w, err)
		return
	}

	stats := struct {
		TotalWorkoutsCompleted int `json:"totalWorkoutsCompleted"`
		WorkoutsThisWeek       int `json:"workoutsThisWeek"`
		TotalMinutesSpent      int `json:"totalMinutesSpent"`
	}{
		TotalWorkoutsCompleted: totalWorkouts,
		WorkoutsThisWeek:       workoutsThisWeek,
		TotalMinutesSpent:      totalSeconds / 60,
	}

	JSONResponse(w, http.StatusOK, stats)
}
