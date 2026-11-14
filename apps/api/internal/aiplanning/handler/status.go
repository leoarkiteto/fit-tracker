package handler

import (
	"net/http"

	"fittracker-api/internal/aiplanning/service"
	"fittracker-api/pkg/httpx"
)

// Status handles GET /api/ai/planning/status.
func Status(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.JSONResponse(w, http.StatusOK, svc.BackendStatus(r.Context()))
	}
}
