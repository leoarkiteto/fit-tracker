package handler

import (
	"net/http"

	"fittracker-api/internal/completedworkout/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetStats handles GET /api/profiles/{profileId}/completed-workouts/stats.
func GetStats(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetStats(r.Context(), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
