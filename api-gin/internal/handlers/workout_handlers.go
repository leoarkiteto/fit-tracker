// Package handlers provides the Gin-based HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"
	"strings"
	"time"

	"fittracker-api-gin/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkoutHandler struct {
	DB *gorm.DB
}

func (h *WorkoutHandler) GetAll(c *gin.Context) {
	profileID := c.Param("profileId")
	var workouts []models.Workout
	if err := h.DB.Preload("Exercises").Where("user_profile_id = ?", profileID).Order("created_at desc").Find(&workouts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	profileID := c.Param("profileId")
	var workout models.Workout
	if err := h.DB.Preload("Exercises").Where("id = ? AND user_profile_id = ?", id, profileID).First(&workout).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) Create(c *gin.Context) {
	profileID := c.Param("profileId")
	var workout models.Workout
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workout.UserProfileID = profileID
	if err := h.DB.Create(&workout).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workout)
}

func (h *WorkoutHandler) Update(c *gin.Context) {
	id := c.Param("id")
	profileID := c.Param("profileId")
	var req models.Workout
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workout models.Workout
	if err := h.DB.Where("id = ? AND user_profile_id = ?", id, profileID).First(&workout).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}

	// In GORM, a common way to update associations is to delete and re-create or use FullSave
	h.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("workout_id = ?", id).Delete(&models.Exercise{})
		workout.Name = req.Name
		workout.Description = req.Description
		workout.Goal = req.Goal
		workout.Days = req.Days
		workout.Exercises = req.Exercises
		return tx.Save(&workout).Error
	})

	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	profileID := c.Param("profileId")
	if err := h.DB.Where("id = ? AND user_profile_id = ?", id, profileID).Delete(&models.Workout{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkoutHandler) GetToday(c *gin.Context) {
	profileID := c.Param("profileId")
	today := strings.ToLower(time.Now().Weekday().String())

	var workouts []models.Workout
	// In GORM we use LIKE for the comma-separated days string
	if err := h.DB.Preload("Exercises").Where("user_profile_id = ? AND days LIKE ?", profileID, "%"+today+"%").Find(&workouts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workouts)
}
