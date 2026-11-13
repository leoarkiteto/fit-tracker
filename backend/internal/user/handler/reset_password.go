package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
)

// resetPasswordSchema is the request body of POST /api/auth/reset-password.
type resetPasswordSchema struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

// ResetPassword handles POST /api/auth/reset-password.
func ResetPassword(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body resetPasswordSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.ResetPassword(r.Context(), body.Email, body.Code, body.NewPassword); err != nil {
			apperror.Write(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
