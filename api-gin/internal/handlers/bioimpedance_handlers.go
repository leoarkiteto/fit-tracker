// Package handlers provides the Gin-based HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"
	"time"

	"fittracker-api-gin/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BioimpedanceHandler struct {
	DB *gorm.DB
}

func (h *BioimpedanceHandler) GetAll(c *gin.Context) {
	profileID := c.Param("profileId")
	var history []models.BioimpedanceData
	if err := h.DB.Where("user_profile_id = ?", profileID).Order("date desc").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *BioimpedanceHandler) Create(c *gin.Context) {
	profileID := c.Param("profileId")
	var data models.BioimpedanceData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.UserProfileID = profileID
	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	if err := h.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, data)
}

func (h *BioimpedanceHandler) GetLatest(c *gin.Context) {
	profileID := c.Param("profileId")
	var data models.BioimpedanceData
	if err := h.DB.Where("user_profile_id = ?", profileID).Order("date desc").First(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *BioimpedanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	profileID := c.Param("profileId")
	if err := h.DB.Where("id = ? AND user_profile_id = ?", id, profileID).Delete(&models.BioimpedanceData{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
