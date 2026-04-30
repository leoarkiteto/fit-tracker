package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"fittracker-api/internal/models"
)

type ProfileHandler struct {
	DB *sql.DB
}

func (h *ProfileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, age, height, current_weight, goal_weight, avatar_url, experience_level, available_days_per_week, preferred_workout_duration, equipment_type FROM UserProfiles`)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var profiles []models.UserProfile
	for rows.Next() {
		var p models.UserProfile
		err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Height, &p.CurrentWeight, &p.GoalWeight, &p.AvatarURL, &p.ExperienceLevel, &p.AvailableDaysPerWeek, &p.PreferredWorkoutDuration, &p.EquipmentType)
		if err != nil {
			serverError(w, err)
			return
		}
		profiles = append(profiles, p)
	}

	JSONResponse(w, http.StatusOK, profiles)
}

func (h *ProfileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.UserProfile
	err := h.DB.QueryRow(`SELECT id, name, age, height, current_weight, goal_weight, avatar_url, experience_level, available_days_per_week, preferred_workout_duration, equipment_type FROM UserProfiles WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Age, &p.Height, &p.CurrentWeight, &p.GoalWeight, &p.AvatarURL, &p.ExperienceLevel, &p.AvailableDaysPerWeek, &p.PreferredWorkoutDuration, &p.EquipmentType)

	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	JSONResponse(w, http.StatusOK, p)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec(`UPDATE UserProfiles SET name = ?, age = ?, height = ?, current_weight = ?, goal_weight = ?, avatar_url = ?, experience_level = ?, available_days_per_week = ?, preferred_workout_duration = ?, equipment_type = ?, updated_at = ? WHERE id = ?`,
		req.Name, req.Age, req.Height, req.CurrentWeight, req.GoalWeight, req.AvatarURL, req.ExperienceLevel, req.AvailableDaysPerWeek, req.PreferredWorkoutDuration, req.EquipmentType, time.Now(), id)

	if err != nil {
		serverError(w, err)
		return
	}

	req.ID = id
	JSONResponse(w, http.StatusOK, req)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	tx, err := h.DB.Begin()
	if err != nil {
		serverError(w, err)
		return
	}

	// Manual cascade: delete all data referencing this profile
	// Exercises are cascaded via Workout delete (ON DELETE CASCADE in schema)
	if _, err = tx.Exec(`DELETE FROM CompletedWorkouts WHERE user_profile_id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}
	if _, err = tx.Exec(`DELETE FROM WaterIntakeEntries WHERE user_profile_id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}
	if _, err = tx.Exec(`DELETE FROM BioimpedanceData WHERE user_profile_id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}
	// Delete workouts (cascades to exercises via ON DELETE CASCADE)
	if _, err = tx.Exec(`DELETE FROM Workouts WHERE user_profile_id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}
	// Unlink user from this profile
	if _, err = tx.Exec(`UPDATE Users SET user_profile_id = NULL WHERE user_profile_id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}
	// Finally delete the profile
	if _, err = tx.Exec(`DELETE FROM UserProfiles WHERE id = ?`, id); err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}

	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
