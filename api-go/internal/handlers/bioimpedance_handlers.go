package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"fittracker-api/internal/models"
	"github.com/google/uuid"
)

type BioimpedanceHandler struct {
	DB *sql.DB
}

func (h *BioimpedanceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")

	rows, err := h.DB.Query(`SELECT id, date, weight, body_fat_percentage, muscle_mass, bone_mass, water_percentage, visceral_fat, bmr, metabolic_age, notes FROM BioimpedanceData WHERE user_profile_id = ? ORDER BY date DESC`, profileID)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var history []models.BioimpedanceData
	for rows.Next() {
		var b models.BioimpedanceData
		err := rows.Scan(&b.ID, &b.Date, &b.Weight, &b.BodyFatPercentage, &b.MuscleMass, &b.BoneMass, &b.WaterPercentage, &b.VisceralFat, &b.BMR, &b.MetabolicAge, &b.Notes)
		if err != nil {
			serverError(w, err)
			return
		}
		b.UserProfileID = profileID
		history = append(history, b)
	}

	JSONResponse(w, http.StatusOK, history)
}

func (h *BioimpedanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	var req models.BioimpedanceData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	_, err := h.DB.Exec(`INSERT INTO BioimpedanceData (id, date, weight, body_fat_percentage, muscle_mass, bone_mass, water_percentage, visceral_fat, bmr, metabolic_age, notes, user_profile_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, req.Date, req.Weight, req.BodyFatPercentage, req.MuscleMass, req.BoneMass, req.WaterPercentage, req.VisceralFat, req.BMR, req.MetabolicAge, req.Notes, profileID)

	if err != nil {
		serverError(w, err)
		return
	}

	req.ID = id
	req.UserProfileID = profileID
	JSONResponse(w, http.StatusCreated, req)
}

func (h *BioimpedanceHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	var b models.BioimpedanceData
	err := h.DB.QueryRow(`SELECT id, date, weight, body_fat_percentage, muscle_mass, bone_mass, water_percentage, visceral_fat, bmr, metabolic_age, notes FROM BioimpedanceData WHERE user_profile_id = ? ORDER BY date DESC LIMIT 1`, profileID).
		Scan(&b.ID, &b.Date, &b.Weight, &b.BodyFatPercentage, &b.MuscleMass, &b.BoneMass, &b.WaterPercentage, &b.VisceralFat, &b.BMR, &b.MetabolicAge, &b.Notes)

	if err != nil {
		http.Error(w, "No data found", http.StatusNotFound)
		return
	}

	b.UserProfileID = profileID
	JSONResponse(w, http.StatusOK, b)
}

func (h *BioimpedanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profileID := r.PathValue("profileId")

	_, err := h.DB.Exec(`DELETE FROM BioimpedanceData WHERE id = ? AND user_profile_id = ?`, id, profileID)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
