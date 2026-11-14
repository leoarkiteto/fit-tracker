package handler

import (
	"net/http"

	"fittracker-api/internal/bioimpedance/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetLatest handles GET /api/profiles/{profileId}/bioimpedance/latest.
func GetLatest(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetLatest(r.Context(), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
