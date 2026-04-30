package main

import (
	"log"

	"fittracker-api-gin/internal/ai"
	"fittracker-api-gin/internal/database"
	"fittracker-api-gin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Database
	db := database.InitDB("fittracker_gin.db")

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

	// 3. Setup Router
	r := gin.Default()

	// Middleware
	r.Use(handlers.CORSMiddleware())

	// Public routes
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
		authGroup.GET("/me", handlers.AuthMiddleware(), authH.Me)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(handlers.AuthMiddleware())
	{
		// Profiles
		api.GET("/profiles", profileH.GetAll)
		api.GET("/profiles/:id", profileH.GetByID)
		api.PUT("/profiles/:id", profileH.Update)
		api.DELETE("/profiles/:id", profileH.Delete)

		// Workouts
		api.GET("/profiles/:profileId/workouts", workoutH.GetAll)
		api.GET("/profiles/:profileId/workouts/today", workoutH.GetToday)
		api.POST("/profiles/:profileId/workouts", workoutH.Create)
		api.GET("/profiles/:profileId/workouts/:id", workoutH.GetByID)
		api.PUT("/profiles/:profileId/workouts/:id", workoutH.Update)
		api.DELETE("/profiles/:profileId/workouts/:id", workoutH.Delete)

		// Completed Workouts
		api.GET("/profiles/:profileId/completed-workouts", completedH.GetAll)
		api.POST("/profiles/:profileId/completed-workouts", completedH.Complete)
		api.GET("/profiles/:profileId/completed-workouts/stats", completedH.GetStats)

		// Bioimpedance
		api.GET("/profiles/:profileId/bioimpedance", bioH.GetAll)
		api.POST("/profiles/:profileId/bioimpedance", bioH.Create)
		api.GET("/profiles/:profileId/bioimpedance/latest", bioH.GetLatest)
		api.DELETE("/profiles/:profileId/bioimpedance/:id", bioH.Delete)

		// Water Intake
		api.GET("/profiles/:profileId/water", waterH.GetDaily)
		api.POST("/profiles/:profileId/water", waterH.Create)
		api.DELETE("/profiles/:profileId/water/:id", waterH.Delete)

		// AI Planning
		api.POST("/ai/planning/generate", aiH.Generate)
		api.POST("/ai/planning/accept", aiH.Accept)
		api.GET("/ai/planning/status", aiH.Status)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// 4. Start Server
	log.Println("Gin/GORM Server starting on :5001...")
	if err := r.Run(":5001"); err != nil {
		log.Fatal(err)
	}
}
