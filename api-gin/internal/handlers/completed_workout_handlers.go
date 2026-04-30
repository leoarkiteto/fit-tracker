// Package handlers provides the Gin-based HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"
	"time"

	"fittracker-api-gin/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompletedWorkoutHandler struct {
	DB *gorm.DB
}

func (h *CompletedWorkoutHandler) GetAll(c *gin.Context) {
	profileID := c.Param("profileId")
	var results []models.CompletedWorkout
	if err := h.DB.Where("user_profile_id = ?", profileID).Order("completed_at desc").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *CompletedWorkoutHandler) Complete(c *gin.Context) {
	profileID := c.Param("profileId")
	var req struct {
		WorkoutID       string `json:"workoutId"`
		DurationSeconds int    `json:"durationSeconds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	completed := models.CompletedWorkout{
		WorkoutID:       req.WorkoutID,
		UserProfileID:   profileID,
		CompletedAt:     time.Now(),
		DurationSeconds: req.DurationSeconds,
	}

	if err := h.DB.Create(&completed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, completed)
}

func (h *CompletedWorkoutHandler) GetStats(c *gin.Context) {
	profileID := c.Param("profileId")

	var totalWorkouts int64
	h.DB.Model(&models.CompletedWorkout{}).Where("user_profile_id = ?", profileID).Count(&totalWorkouts)

	var workoutsThisWeek int64
	lastWeek := time.Now().AddDate(0, 0, -7)
	h.DB.Model(&models.CompletedWorkout{}).Where("user_profile_id = ? AND completed_at >= ?", profileID, lastWeek).Count(&workoutsThisWeek)

	var totalSeconds int64
	h.DB.Model(&models.CompletedWorkout{}).Where("user_profile_id = ?", profileID).Select("COALESCE(SUM(duration_seconds), 0)").Row().Scan(&totalSeconds)

	c.JSON(http.StatusOK, gin.H{
		"totalWorkoutsCompleted": totalWorkouts,
		"workoutsThisWeek":       workoutsThisWeek,
		"totalMinutesSpent":      totalSeconds / 60,
	})
}
