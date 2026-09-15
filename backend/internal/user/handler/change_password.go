package handler

import (
	"encoding/json"
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// changePasswordSchema is the request body of PATCH /api/auth/change-password.
type changePasswordSchema struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// ChangePassword handles PATCH /api/auth/change-password.
func ChangePassword(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := httpx.UserID(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body changePasswordSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := svc.ChangePassword(r.Context(), userID, body.CurrentPassword, body.NewPassword)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
	}
}
