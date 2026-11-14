package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// forgotPasswordSchema is the request body of POST /api/auth/forgot-password.
type forgotPasswordSchema struct {
	Email string `json:"email"`
}

// ForgotPassword handles POST /api/auth/forgot-password. The response is the
// same whether or not an account exists for the email.
func ForgotPassword(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body forgotPasswordSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.RequestPasswordReset(r.Context(), body.Email); err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, map[string]string{
			"message": "Se houver uma conta com este e-mail, enviamos um código",
		})
	}
}
