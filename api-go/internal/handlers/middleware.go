package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"fittracker-api/internal/auth"
)

type contextKey string

const (
	UserIDKey          contextKey = "userID"
	ProfileIDKey       contextKey = "profileID"
	maxRequestBodySize           = 10 << 20 // 10 MB limit for request bodies
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Safe type assertions with checks
		sub, ok := (*claims)["sub"].(string)
		if !ok || sub == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, sub)

		profileID, _ := (*claims)["profileId"].(string)
		ctx = context.WithValue(ctx, ProfileIDKey, profileID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// JSONResponse writes a JSON response with the given status code and data.
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// LimitRequestBody wraps a handler with a request body size limit (10 MB).
func LimitRequestBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		next.ServeHTTP(w, r)
	}
}

// serverError logs the real error and returns a generic 500 response.
func serverError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
