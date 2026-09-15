package handler

import (
	"net/http"

	"fittracker-api/internal/water/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// GetDaily handles GET /api/profiles/{profileId}/water?date=YYYY-MM-DD.
func GetDaily(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetDaily(
			r.Context(),
			r.PathValue("profileId"),
			r.URL.Query().Get("date"),
		)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
