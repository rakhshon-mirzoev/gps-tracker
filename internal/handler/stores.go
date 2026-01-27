package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) CreateStore(c *gin.Context) {
	var dto *model.CreateStore
	user_id, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("CreateStore: failed to check staff: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("CreateStore: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.s.Store.Create(c.Request.Context(), dto, user_id)
	if err != nil {
		logger.Log.Errorf("CreateStore: failed to create store: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "create complete"})
}
