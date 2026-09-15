package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/aiplanning"
	"fittracker-api/internal/aiplanning/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// acceptSchema is the request body of POST /api/ai/planning/accept. PlanID is
// accepted for compatibility and not used, as before.
type acceptSchema struct {
	PlanID   string                      `json:"planId"`
	Workouts []aiplanning.PlannedWorkout `json:"workouts"`
}

// Accept handles POST /api/ai/planning/accept. The target profile comes from
// the token claims rather than the URL.
func Accept(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := httpx.ProfileID(r)
		if profileID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body acceptSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.Accept(r.Context(), profileID, body.Workouts); err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(
			w,
			http.StatusOK,
			map[string]string{"message": "Plan accepted and workouts created"},
		)
	}
}
