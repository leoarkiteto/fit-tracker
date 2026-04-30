package handlers

import (
	"database/sql"
	"encoding/json"
	"fittracker-api/internal/auth"
	"fittracker-api/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	profileID := uuid.New().String()
	userID := uuid.New().String()

	tx, err := h.DB.Begin()
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = tx.Exec(`INSERT INTO UserProfiles (id, name) VALUES (?, ?)`, profileID, req.Name)
	if err != nil {
		tx.Rollback()
		serverError(w, err)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO Users (id, email, password_hash, name, user_profile_id) VALUES (?, ?, ?, ?, ?)`,
		userID,
		req.Email,
		hashedPassword,
		req.Name,
		profileID,
	)
	if err != nil {
		tx.Rollback()
		// Use generic message to avoid email enumeration
		http.Error(w, "Registration failed", http.StatusConflict)
		return
	}

	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	user := models.User{
		ID:            userID,
		Email:         req.Email,
		Name:          req.Name,
		UserProfileID: &profileID,
	}

	token, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, http.StatusCreated, models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserDto{
			ID:        userID,
			Email:     req.Email,
			Name:      req.Name,
			ProfileID: &profileID,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	var profileID sql.NullString
	err := h.DB.QueryRow(`SELECT id, email, password_hash, name, user_profile_id FROM Users WHERE email = ?`, req.Email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &profileID)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if profileID.Valid {
		user.UserProfileID = &profileID.String
	}

	_, _ = h.DB.Exec(`UPDATE Users SET last_login_at = ? WHERE id = ?`, time.Now(), user.ID)

	token, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, http.StatusOK, models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserDto{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			ProfileID: user.UserProfileID,
		},
	})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	// Fetch current password hash
	var currentHash string
	err := h.DB.QueryRow(`SELECT password_hash FROM Users WHERE id = ?`, userID).Scan(&currentHash)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !auth.CheckPasswordHash(req.CurrentPassword, currentHash) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash and update
	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		serverError(w, err)
		return
	}

	_, err = h.DB.Exec(`UPDATE Users SET password_hash = ? WHERE id = ?`, newHash, userID)
	if err != nil {
		serverError(w, err)
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var user models.UserDto
	var profileID sql.NullString
	err := h.DB.QueryRow(`SELECT id, email, name, user_profile_id FROM Users WHERE id = ?`, userID).
		Scan(&user.ID, &user.Email, &user.Name, &profileID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if profileID.Valid {
		user.ProfileID = &profileID.String
	}

	JSONResponse(w, http.StatusOK, user)
}
