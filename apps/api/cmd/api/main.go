package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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
	"fittracker-api/internal/user"
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

// splitList splits a comma-separated value, trimming each item and dropping empties.
func splitList(value string) []string {
	var items []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, item)
		}
	}
	return items
}

func main() {
	// 1. Infrastructure: the database and the token signer (JWT_SECRET is
	// required for the server to start).
	db := database.InitDB("fittracker.db")
	defer db.Close()

	secret := os.Getenv("JWT_SECRET")
	signer, err := auth.NewSigner(secret)
	if err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}
	resetKey := auth.DeriveKey(secret, "password-reset")

	log.Println("WARNING: development mailer active — password reset codes are written to the server log. Do not deploy as-is.")

	// 2. Outbound adapters.
	ollamaClient, err := aiinfra.NewOllamaClient("http://localhost:11434", "llama3.2")
	if err != nil {
		log.Fatalf("Failed to initialize AI Client: %v", err)
	}

	// 3. Router: every slice registers its own routes.
	mux := http.NewServeMux()

	// User slice
	userSvc := userservice.New(
		userinfra.NewUserRepository(db),
		userinfra.NewPasswordResetRepository(db),
		userinfra.NewLogMailer(),
		signer,
		resetKey,
	)

	// Social login providers are optional: each stays a nil interface (and its
	// endpoint answers 503) unless fully configured.
	var google user.GoogleVerifier
	if ids := splitList(os.Getenv("GOOGLE_CLIENT_IDS")); len(ids) > 0 {
		google = userinfra.NewGoogleVerifier(ids)
		log.Println("Social login: Google enabled")
	} else {
		log.Println("Social login: Google disabled (GOOGLE_CLIENT_IDS not set)")
	}

	var github user.GitHubExchanger
	githubID := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID"))
	githubSecret := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET"))
	if githubID != "" && githubSecret != "" {
		github = userinfra.NewGitHubExchanger(githubID, githubSecret)
		log.Println("Social login: GitHub enabled")
	} else {
		log.Println("Social login: GitHub disabled (GITHUB_CLIENT_ID/GITHUB_CLIENT_SECRET not set)")
	}

	userhandler.Router(mux, signer, userSvc.WithSocialLogin(userinfra.NewIdentityRepository(db), google, github))

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
