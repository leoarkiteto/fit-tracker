package handler

import (
	"net/http"

	"fittracker-api/internal/workout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetToday handles GET /api/profiles/{profileId}/workouts/today.
func GetToday(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetToday(r.Context(), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
