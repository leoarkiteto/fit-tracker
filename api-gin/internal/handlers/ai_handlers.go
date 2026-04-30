// Package handlers provides the Gin-based HTTP request handlers for the API endpoints.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fittracker-api-gin/internal/ai"
	"fittracker-api-gin/internal/models"
)

type AIHandler struct {
	DB       *gorm.DB
	AIClient *ai.LangChainClient
}

const systemPrompt = `You are an experienced fitness coach. Respond ONLY with valid JSON.`

const promptTemplate = `Based on user profile, create a personalized weekly workout plan.
User: %s, Age: %v, Weight: %vkg, Goal: %s.
Generate JSON structure with: summary, rationale, workouts (name, description, goal, days, exercises).`

func (h *AIHandler) Generate(c *gin.Context) {
	var req struct {
		UserProfileID string `json:"userProfileId"`
		Goal          string `json:"goal"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var p models.UserProfile
	if err := h.DB.First(&p, "id = ?", req.UserProfileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	weight := 0.0
	if p.CurrentWeight != nil {
		weight = *p.CurrentWeight
	}
	age := 0
	if p.Age != nil {
		age = *p.Age
	}

	userPrompt := fmt.Sprintf(promptTemplate, p.Name, age, weight, req.Goal)

	responseText, err := h.AIClient.GenerateWorkoutPlan(
		c.Request.Context(),
		systemPrompt,
		userPrompt,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "AI Generation failed: " + err.Error()},
		)
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(responseText))
}

func (h *AIHandler) Accept(c *gin.Context) {
	var req struct {
		PlanID        string           `json:"planId"`
		UserProfileID string           `json:"userProfileId"`
		Workouts      []models.Workout `json:"workouts"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		for _, workout := range req.Workouts {
			workout.UserProfileID = req.UserProfileID
			// In GORM, just creating the workout will create exercises if they are in the struct
			if err := tx.Create(&workout).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan accepted and workouts created"})
}

func (h *AIHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"available": true,
		"provider":  "Ollama (via LangChainGo) + Gin + GORM",
	})
}
