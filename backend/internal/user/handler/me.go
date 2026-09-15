package handler

import (
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// Me handles GET /api/auth/me.
func Me(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.Me(r.Context(), httpx.UserID(r))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
