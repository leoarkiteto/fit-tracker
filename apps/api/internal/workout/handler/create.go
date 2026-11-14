package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/workout"
	"fittracker-api/internal/workout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// Create handles POST /api/profiles/{profileId}/workouts. As before, the
// response echoes the submitted body with the generated ids.
func Create(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body workout.Workout
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.Create(r.Context(), r.PathValue("profileId"), &body); err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusCreated, body)
	}
}
