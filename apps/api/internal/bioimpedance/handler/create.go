package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/bioimpedance"
	"fittracker-api/internal/bioimpedance/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// Create handles POST /api/profiles/{profileId}/bioimpedance. As before, the
// response echoes the submitted body with the generated id, the normalised date
// and the profile id.
func Create(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body bioimpedance.Measurement
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Create(r.Context(), r.PathValue("profileId"), &body)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusCreated, result)
	}
}
