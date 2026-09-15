package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"fittracker-api/internal/water/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// createSchema is the request body of POST /api/profiles/{profileId}/water.
type createSchema struct {
	AmountMl   int        `json:"amountMl"`
	ConsumedAt *time.Time `json:"consumedAt"`
}

// Create handles POST /api/profiles/{profileId}/water.
func Create(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body createSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Create(
			r.Context(),
			r.PathValue("profileId"),
			body.AmountMl,
			body.ConsumedAt,
		)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusCreated, result)
	}
}
