package handler

import (
	"net/http"

	"fittracker-api/internal/bioimpedance/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetAll handles GET /api/profiles/{profileId}/bioimpedance.
func GetAll(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetAll(r.Context(), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
