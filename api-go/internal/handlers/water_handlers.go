// Package handlers with
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"fittracker-api/internal/models"
)

type WaterHandler struct {
	DB *sql.DB
}

func (h *WaterHandler) GetDaily(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	dateStr := r.URL.Query().Get("date")

	var date time.Time
	var err error
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
	} else {
		date = time.Now()
	}
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	rows, err := h.DB.Query(
		`SELECT id, amount_ml, consumed_at, note FROM WaterIntakeEntries WHERE user_profile_id = ? AND consumed_at >= ? AND consumed_at < ?`,
		profileID,
		start,
		end,
	)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()

	var entries []models.WaterIntakeEntry
	totalMl := 0
	for rows.Next() {
		var e models.WaterIntakeEntry
		err := rows.Scan(&e.ID, &e.AmountMl, &e.ConsumedAt, &e.Note)
		if err != nil {
			serverError(w, err)
			return
		}
		totalMl += e.AmountMl
		entries = append(entries, e)
	}

	// Calculate goal based on weight (35ml per kg)
	var weight float64
	err = h.DB.QueryRow(`SELECT COALESCE(current_weight, 70.0) FROM UserProfiles WHERE id = ?`, profileID).
		Scan(&weight)
	if err != nil {
		weight = 70.0
	}
	goalMl := int(weight * 35)

	summary := struct {
		Date    string                    `json:"date"`
		TotalMl int                       `json:"totalMl"`
		GoalMl  int                       `json:"goalMl"`
		Entries []models.WaterIntakeEntry `json:"entries"`
	}{
		Date:    start.Format("2006-01-02"),
		TotalMl: totalMl,
		GoalMl:  goalMl,
		Entries: entries,
	}

	JSONResponse(w, http.StatusOK, summary)
}

func (h *WaterHandler) Create(w http.ResponseWriter, r *http.Request) {
	profileID := r.PathValue("profileId")
	var req struct {
		AmountMl   int        `json:"amountMl"`
		ConsumedAt *time.Time `json:"consumedAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	consumedAt := time.Now()
	if req.ConsumedAt != nil {
		consumedAt = *req.ConsumedAt
	}

	_, err := h.DB.Exec(
		`INSERT INTO WaterIntakeEntries (id, user_profile_id, amount_ml, consumed_at) VALUES (?, ?, ?, ?)`,
		id,
		profileID,
		req.AmountMl,
		consumedAt,
	)
	if err != nil {
		serverError(w, err)
		return
	}

	JSONResponse(w, http.StatusCreated, models.WaterIntakeEntry{
		ID:            id,
		UserProfileID: profileID,
		AmountMl:      req.AmountMl,
		ConsumedAt:    consumedAt,
	})
}

func (h *WaterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profileID := r.PathValue("profileId")

	_, err := h.DB.Exec(
		`DELETE FROM WaterIntakeEntries WHERE id = ? AND user_profile_id = ?`,
		id,
		profileID,
	)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
