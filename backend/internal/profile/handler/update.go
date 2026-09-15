package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/profile"
	"fittracker-api/internal/profile/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// Update handles PUT /api/profiles/{id}. As before, the response echoes the
// submitted body with the id from the path.
func Update(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body profile.Profile
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.ID = r.PathValue("id")

		if err := svc.Update(r.Context(), &body); err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, body)
	}
}
