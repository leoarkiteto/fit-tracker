package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/aiplanning/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// generateSchema is the request body of POST /api/ai/planning/generate.
type generateSchema struct {
	UserProfileID string `json:"userProfileId"`
	Goal          string `json:"goal"`
}

// Generate handles POST /api/ai/planning/generate.
func Generate(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body generateSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Generate(r.Context(), body.UserProfileID, body.Goal)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
