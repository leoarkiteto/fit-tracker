package main

import (
	"log"
	"net/http"

	"fittracker-api/internal/ai"
	"fittracker-api/internal/database"
	"fittracker-api/internal/handlers"
)

func main() {
	// 1. Initialize Database
	db := database.InitDB("fittracker.db")
	defer db.Close()

	// 2. Initialize Clients & Handlers
	aiClient, err := ai.NewLangChainClient("http://localhost:11434", "llama3.2")
	if err != nil {
		log.Fatalf("Failed to initialize AI Client: %v", err)
	}
	
	authH := &handlers.AuthHandler{DB: db}
	profileH := &handlers.ProfileHandler{DB: db}
	workoutH := &handlers.WorkoutHandler{DB: db}
	completedH := &handlers.CompletedWorkoutHandler{DB: db}
	bioH := &handlers.BioimpedanceHandler{DB: db}
	waterH := &handlers.WaterHandler{DB: db}
	aiH := &handlers.AIHandler{DB: db, AIClient: aiClient}

	// 3. Setup Router (Go 1.22+ syntax)
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", handlers.LimitRequestBody(authH.Register))
	mux.HandleFunc("POST /api/auth/login", handlers.LimitRequestBody(authH.Login))
	mux.HandleFunc("GET /api/auth/me", handlers.AuthMiddleware(authH.Me))
	mux.HandleFunc("PATCH /api/auth/change-password", handlers.LimitRequestBody(handlers.AuthMiddleware(authH.ChangePassword)))

	// Profile routes
	mux.HandleFunc("GET /api/profiles", handlers.AuthMiddleware(profileH.GetAll))
	mux.HandleFunc("GET /api/profiles/{id}", handlers.AuthMiddleware(profileH.GetByID))
	mux.HandleFunc("PUT /api/profiles/{id}", handlers.LimitRequestBody(handlers.AuthMiddleware(profileH.Update)))
	mux.HandleFunc("DELETE /api/profiles/{id}", handlers.AuthMiddleware(profileH.Delete))

	// Workout routes
	mux.HandleFunc("GET /api/profiles/{profileId}/workouts", handlers.AuthMiddleware(workoutH.GetAll))
	mux.HandleFunc("GET /api/profiles/{profileId}/workouts/today", handlers.AuthMiddleware(workoutH.GetToday))
	mux.HandleFunc("POST /api/profiles/{profileId}/workouts", handlers.LimitRequestBody(handlers.AuthMiddleware(workoutH.Create)))
	mux.HandleFunc("GET /api/profiles/{profileId}/workouts/{id}", handlers.AuthMiddleware(workoutH.GetByID))
	mux.HandleFunc("PUT /api/profiles/{profileId}/workouts/{id}", handlers.LimitRequestBody(handlers.AuthMiddleware(workoutH.Update)))
	mux.HandleFunc("DELETE /api/profiles/{profileId}/workouts/{id}", handlers.AuthMiddleware(workoutH.Delete))

	// Completed Workout routes
	mux.HandleFunc("GET /api/profiles/{profileId}/completed-workouts", handlers.AuthMiddleware(completedH.GetAll))
	mux.HandleFunc("POST /api/profiles/{profileId}/completed-workouts", handlers.LimitRequestBody(handlers.AuthMiddleware(completedH.Complete)))
	mux.HandleFunc("GET /api/profiles/{profileId}/completed-workouts/stats", handlers.AuthMiddleware(completedH.GetStats))

	// Bioimpedance routes
	mux.HandleFunc("GET /api/profiles/{profileId}/bioimpedance", handlers.AuthMiddleware(bioH.GetAll))
	mux.HandleFunc("POST /api/profiles/{profileId}/bioimpedance", handlers.LimitRequestBody(handlers.AuthMiddleware(bioH.Create)))
	mux.HandleFunc("GET /api/profiles/{profileId}/bioimpedance/latest", handlers.AuthMiddleware(bioH.GetLatest))
	mux.HandleFunc("DELETE /api/profiles/{profileId}/bioimpedance/{id}", handlers.AuthMiddleware(bioH.Delete))

	// Water Intake routes
	mux.HandleFunc("GET /api/profiles/{profileId}/water", handlers.AuthMiddleware(waterH.GetDaily))
	mux.HandleFunc("POST /api/profiles/{profileId}/water", handlers.LimitRequestBody(handlers.AuthMiddleware(waterH.Create)))
	mux.HandleFunc("DELETE /api/profiles/{profileId}/water/{id}", handlers.AuthMiddleware(waterH.Delete))

	// AI Planning routes
	mux.HandleFunc("POST /api/ai/planning/generate", handlers.LimitRequestBody(handlers.AuthMiddleware(aiH.Generate)))
	mux.HandleFunc("POST /api/ai/planning/accept", handlers.LimitRequestBody(handlers.AuthMiddleware(aiH.Accept)))
	mux.HandleFunc("GET /api/ai/planning/status", aiH.Status)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 4. Wrap with CORS Middleware
	handler := corsMiddleware(mux)

	// 5. Start Server
	log.Println("Go Server starting on :5000...")
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
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
