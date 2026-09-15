package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/workout"
	"fittracker-api/internal/workout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// Update handles PUT /api/profiles/{profileId}/workouts/{id}. As before, the
// response echoes the submitted body with the id from the path.
func Update(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body workout.Workout
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := svc.Update(r.Context(), r.PathValue("id"), r.PathValue("profileId"), &body)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, body)
	}
}
