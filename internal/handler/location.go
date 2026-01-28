package handler

import (
	"net/http"

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

	err = h.s.Location.Create(c.Request.Context(), loc, user_id)
	if err != nil {
		logger.Log.Errorf("Location: failed to create location: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "location created: " + loc.ID.String()})
}
