// Package handlers provides the Gin-based HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fittracker-api-gin/internal/models"
)

type WaterHandler struct {
	DB *gorm.DB
}

func (h *WaterHandler) GetDaily(c *gin.Context) {
	profileID := c.Param("profileId")
	dateStr := c.Query("date")

	var date time.Time
	var err error
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
	} else {
		date = time.Now()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	var entries []models.WaterIntakeEntry
	if err := h.DB.Where("user_profile_id = ? AND consumed_at >= ? AND consumed_at < ?", profileID, start, end).Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalMl := 0
	for _, e := range entries {
		totalMl += e.AmountMl
	}

	var profile models.UserProfile
	var weight float64 = 70.0
	if err := h.DB.Select("current_weight").First(&profile, "id = ?", profileID).Error; err == nil &&
		profile.CurrentWeight != nil {
		weight = *profile.CurrentWeight
	}
	goalMl := int(weight * 35)

	c.JSON(http.StatusOK, gin.H{
		"date":    start.Format("2006-01-02"),
		"totalMl": totalMl,
		"goalMl":  goalMl,
		"entries": entries,
	})
}

func (h *WaterHandler) Create(c *gin.Context) {
	profileID := c.Param("profileId")
	var req struct {
		AmountMl   int        `json:"amountMl"`
		ConsumedAt *time.Time `json:"consumedAt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumedAt := time.Now()
	if req.ConsumedAt != nil {
		consumedAt = *req.ConsumedAt
	}

	entry := models.WaterIntakeEntry{
		UserProfileID: profileID,
		AmountMl:      req.AmountMl,
		ConsumedAt:    consumedAt,
	}

	if err := h.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *WaterHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	profileID := c.Param("profileId")
	if err := h.DB.Where("id = ? AND user_profile_id = ?", id, profileID).Delete(&models.WaterIntakeEntry{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
