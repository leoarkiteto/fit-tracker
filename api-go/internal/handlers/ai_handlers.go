package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"fittracker-api/internal/ai"
	"fittracker-api/internal/models"
	"github.com/google/uuid"
)

const aiTimeout = 60 * time.Second

type AIHandler struct {
	DB     *sql.DB
	AIClient *ai.LangChainClient
}

const systemPrompt = `You are an experienced fitness coach. Respond ONLY with valid JSON.`

const promptTemplate = `Based on user profile, create a personalized weekly workout plan.
User: %s, Age: %v, Weight: %vkg, Goal: %s.
Generate JSON structure with: summary, rationale, workouts (name, description, goal, days, exercises).`

func (h *AIHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserProfileID string `json:"userProfileId"`
		Goal          string `json:"goal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p models.UserProfile
	err := h.DB.QueryRow(`SELECT name, age, current_weight FROM UserProfiles WHERE id = ?`, req.UserProfileID).
		Scan(&p.Name, &p.Age, &p.CurrentWeight)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	userPrompt := fmt.Sprintf(promptTemplate, p.Name, p.Age, p.CurrentWeight, req.Goal)
	
	// Timeout for AI generation to prevent hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), aiTimeout)
	defer cancel()

	responseText, err := h.AIClient.GenerateWorkoutPlan(ctx, systemPrompt, userPrompt)
	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responseText))
}

func (h *AIHandler) Accept(w http.ResponseWriter, r *http.Request) {
	// Extract profileId from JWT claims (set by AuthMiddleware)
	profileID, ok := r.Context().Value(ProfileIDKey).(string)
	if !ok || profileID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		PlanID   string           `json:"planId"`
		Workouts []models.Workout `json:"workouts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		serverError(w, err)
		return
	}

	for _, workout := range req.Workouts {
		workoutID := uuid.New().String()
		daysStr := strings.Join(workout.Days, ",")
		_, err = tx.Exec(`INSERT INTO Workouts (id, name, description, goal, days, user_profile_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			workoutID, workout.Name, workout.Description, workout.Goal, daysStr, profileID, time.Now(), time.Now())
		if err != nil {
			tx.Rollback()
			serverError(w, err)
			return
		}

		for _, ex := range workout.Exercises {
			exID := uuid.New().String()
			_, err = tx.Exec(`INSERT INTO Exercises (id, name, muscle_group, sets, reps, weight, rest_seconds, notes, workout_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				exID, ex.Name, ex.MuscleGroup, ex.Sets, ex.Reps, ex.Weight, ex.RestSeconds, ex.Notes, workoutID)
			if err != nil {
				tx.Rollback()
				serverError(w, err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Plan accepted and workouts created"})
}

func (h *AIHandler) Status(w http.ResponseWriter, r *http.Request) {
	// Status could be more sophisticated with LangChain
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"available": true,
		"provider":  "Ollama (via LangChainGo)",
	})
}
