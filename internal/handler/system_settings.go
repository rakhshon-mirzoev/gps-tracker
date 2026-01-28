package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) CreateSystemSetting(c *gin.Context) {
	var ss *model.SystemSettings

	if err := c.ShouldBindJSON(&ss); err != nil {
		logger.Log.Errorf("CreateSystemSetting: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.SystemSettings.Create(c.Request.Context(), ss); err != nil {
		logger.Log.Errorf("CreateSystemSetting: failed to create: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + ss.ID.String()})

}

func (h *Handler) ListSystemSettings(c *gin.Context) {

	ss, err := h.s.SystemSettings.List(c.Request.Context())
	if err != nil {
		logger.Log.Errorf("ListSystemSettings: failed to list: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"system_settings": ss})

}

func (h *Handler) DeleteSystemSetting(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.Errorf("DeleteSystemSetting: failed to parse id: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.SystemSettings.Delete(c.Request.Context(), parsedID); err != nil {
		logger.Log.Errorf("DeleteSystemSetting: failed to delete: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateSystemSetting(c *gin.Context) {
	var dto model.UpdateSystemSettings

	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.Errorf("UpdateSystemSetting: failed to parse id: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("UpdateSystemSetting: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.SystemSettings.Update(c.Request.Context(), parsedID, dto); err != nil {
		logger.Log.Errorf("UpdateSystemSetting: failed to update: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
