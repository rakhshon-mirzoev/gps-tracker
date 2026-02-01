package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) ListShedules(c *gin.Context) {
	user_id, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("ListShedules: failed to check user: %w", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "error while user check"})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		logger.Log.Errorf("ListShedules: failed to getLimitOffset: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	schedules, total, err := h.s.Schedule.ListByStaffID(c.Request.Context(), user_id, lim, off)
	if err != nil {
		logger.Log.Errorf("ListShedules: failed to ListByStaffID: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules, "total": total})

}

func (h *Handler) AssignSchedule(c *gin.Context) {
	var dto *model.CreateSchedule

	adminID, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("AssignSchedule: failed bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.s.Schedule.Create(c.Request.Context(), dto, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Errorf("AssignSchedule: failed create: %w", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": dto})
}
