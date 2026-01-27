package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) VisitCheckIn(c *gin.Context) {
	var visit = &model.CreateVisit{}
	user_id, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("failed to check staff user: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("failed to check bind visit: %w", err)
		return
	}

	err = h.s.Visit.Create(c.Request.Context(), visit, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("failed to create visit: %w", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "check-in completed"})

}

func (h *Handler) VisitPhoto(c *gin.Context) {
	var dto model.UpdateVisit

	vID := c.Param("id")

	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении файлов: " + err.Error()})
		logger.Log.Errorf("multipart form: %w", err)
		return
	}

	err = h.s.Visit.VisitPhoto(c.Request.Context(), vID, dto, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("failed to update visit: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update complete"})

}

func (h *Handler) VisitComplete(c *gin.Context) {
	var dto model.UpdateVisit

	vID := c.Param("id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("failed to bind dto: %w", err)
		return
	}

	err := h.s.Visit.VisitComplete(c.Request.Context(), vID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("failed to update visit: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "visit complete"})
}
