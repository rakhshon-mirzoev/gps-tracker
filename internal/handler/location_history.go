package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) Location(c *gin.Context) {
	var loc *model.LocationHistory
	user_id, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("Location: failed to check user: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&loc); err != nil {
		logger.Log.Errorf("Location: failed to bind location: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	loc.StaffID = user_id

	err = h.s.Location.Create(c.Request.Context(), loc, user_id)
	if err != nil {
		logger.Log.Errorf("Location: failed to create location: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "location created: " + loc.ID.String()})
}

func (h *Handler) AllStaffOnlineOffline(c *gin.Context) {
	id, err := checkSystemUser(c)
	if err != nil {
		logger.Log.Errorf("AllStaffOnlineOffline: failed to check System User: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check System User"})
		return
	}

	beginStr := c.Query("begin")
	endStr := c.Query("end")

	begin, err := time.Parse("2006-01-02", beginStr)
	if err != nil {
		logger.Log.Errorf("AllStaffOnlineOffline: failed to parse begin date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid begin date"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		logger.Log.Errorf("AllStaffOnlineOffline: failed to parse end date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	data, err := h.s.Location.GetAllStaffOnlineOffline(c.Request.Context(), id, begin, end)
	if err != nil {
		logger.Log.Errorf("AllStaffOnlineOffline: failed to get durations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get staff durations"})
		return
	}

	result := make([]gin.H, 0, len(data))
	for staffID, durations := range data {
		result = append(result, gin.H{
			"staff_id": staffID,
			"online":   durations["online"].Minutes(),
			"offline":  durations["offline"].Minutes(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
