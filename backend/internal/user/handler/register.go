package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// registerSchema is the request body of POST /api/auth/register.
type registerSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Register handles POST /api/auth/register.
func Register(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body registerSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Register(r.Context(), body.Email, body.Password, body.Name)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusCreated, result)
	}
}
