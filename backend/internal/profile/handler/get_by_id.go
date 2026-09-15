package handler

import (
	"net/http"

	"fittracker-api/internal/profile/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetByID handles GET /api/profiles/{id}.
func GetByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetByID(r.Context(), r.PathValue("id"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
