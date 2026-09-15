package main

import (
	"log"
	"net/http"
	"os"

	aihandler "fittracker-api/internal/aiplanning/handler"
	aiinfra "fittracker-api/internal/aiplanning/infrastructure"
	aiservice "fittracker-api/internal/aiplanning/service"
	biohandler "fittracker-api/internal/bioimpedance/handler"
	bioinfra "fittracker-api/internal/bioimpedance/infrastructure"
	bioservice "fittracker-api/internal/bioimpedance/service"
	completedhandler "fittracker-api/internal/completedworkout/handler"
	completedinfra "fittracker-api/internal/completedworkout/infrastructure"
	completedservice "fittracker-api/internal/completedworkout/service"
	profilehandler "fittracker-api/internal/profile/handler"
	profileinfra "fittracker-api/internal/profile/infrastructure"
	profileservice "fittracker-api/internal/profile/service"
	userhandler "fittracker-api/internal/user/handler"
	userinfra "fittracker-api/internal/user/infrastructure"
	userservice "fittracker-api/internal/user/service"
	waterhandler "fittracker-api/internal/water/handler"
	waterinfra "fittracker-api/internal/water/infrastructure"
	waterservice "fittracker-api/internal/water/service"
	workouthandler "fittracker-api/internal/workout/handler"
	workoutinfra "fittracker-api/internal/workout/infrastructure"
	workoutservice "fittracker-api/internal/workout/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/database"
	"fittracker-api/pkg/httpx"
)

func main() {
	// 1. Infrastructure: the database and the token signer (JWT_SECRET is
	// required for the server to start).
	db := database.InitDB("fittracker.db")
	defer db.Close()

	signer, err := auth.NewSigner(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	// 2. Outbound adapters.
	ollamaClient, err := aiinfra.NewOllamaClient("http://localhost:11434", "llama3.2")
	if err != nil {
		log.Fatalf("Failed to initialize AI Client: %v", err)
	}

	// 3. Router: every slice registers its own routes.
	mux := http.NewServeMux()

	// User slice
	userhandler.Router(mux, signer, userservice.New(userinfra.NewUserRepository(db), signer))

	// Profile slice
	profilehandler.Router(mux, signer, profileservice.New(profileinfra.NewProfileRepository(db)))

	// Workout slice
	workouthandler.Router(mux, signer, workoutservice.New(workoutinfra.NewWorkoutRepository(db)))

	// Completed workout slice
	completedhandler.Router(
		mux,
		signer,
		completedservice.New(completedinfra.NewCompletedWorkoutRepository(db)),
	)

	// Bioimpedance slice
	biohandler.Router(mux, signer, bioservice.New(bioinfra.NewMeasurementRepository(db)))

	// Water intake slice
	waterhandler.Router(mux, signer, waterservice.New(waterinfra.NewEntryRepository(db)))

	// AI planning slice
	aihandler.Router(mux, signer, aiservice.New(
		aiinfra.NewProfileReader(db),
		aiinfra.NewWorkoutWriter(db),
		ollamaClient,
	))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 4. Start server
	log.Println("Go Server starting on :5000...")
	if err := http.ListenAndServe(":5000", httpx.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}
