package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/completedworkout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// completeSchema is the request body of POST /api/profiles/{profileId}/completed-workouts.
type completeSchema struct {
	WorkoutID       string `json:"workoutId"`
	DurationSeconds int    `json:"durationSeconds"`
}

// Complete handles POST /api/profiles/{profileId}/completed-workouts.
func Complete(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body completeSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Complete(
			r.Context(),
			r.PathValue("profileId"),
			body.WorkoutID,
			body.DurationSeconds,
		)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusCreated, result)
	}
}
