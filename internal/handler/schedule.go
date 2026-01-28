package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
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
