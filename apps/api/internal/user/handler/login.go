package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// loginSchema is the request body of POST /api/auth/login.
type loginSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /api/auth/login.
func Login(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body loginSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := svc.Login(r.Context(), body.Email, body.Password)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
