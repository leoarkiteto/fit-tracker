package handler

import (
	"net/http"

	"fittracker-api/internal/workout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetByID handles GET /api/profiles/{profileId}/workouts/{id}.
func GetByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetByID(r.Context(), r.PathValue("id"), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
