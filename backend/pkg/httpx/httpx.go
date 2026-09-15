// Package httpx holds the transport helpers every slice uses: request context
// keys, the authentication guard, body limits, JSON rendering and CORS.
package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"fittracker-api/pkg/auth"
)

// ContextKey is the type of the keys this package stores in a request context.
type ContextKey string

const (
	// ContextUserID holds the authenticated user id.
	ContextUserID ContextKey = "userID"
	// ContextProfileID holds the profile id carried by the token.
	ContextProfileID ContextKey = "profileID"

	maxRequestBodySize = 10 << 20 // 10 MB limit for request bodies
)

// UserID returns the authenticated user id, or "" when there is none.
func UserID(r *http.Request) string {
	id, _ := r.Context().Value(ContextUserID).(string)
	return id
}

// ProfileID returns the profile id carried by the token, or "" when there is none.
func ProfileID(r *http.Request) string {
	id, _ := r.Context().Value(ContextProfileID).(string)
	return id
}

// RequireAuth rejects requests without a valid bearer token and injects the
// token claims into the request context.
func RequireAuth(signer *auth.Signer, next http.HandlerFunc) http.HandlerFunc {
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

		claims, err := signer.ValidateToken(parts[1])
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
		ctx := context.WithValue(r.Context(), ContextUserID, sub)

		profileID, _ := (*claims)["profileId"].(string)
		ctx = context.WithValue(ctx, ContextProfileID, profileID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// JSONResponse writes a JSON response with the given status code and data.
func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// LimitBody wraps a handler with the maximum request body size (10 MB).
func LimitBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		next.ServeHTTP(w, r)
	}
}

// CORS allows the mobile client to call the API from any origin.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
