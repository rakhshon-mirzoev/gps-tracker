package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
)

func (h *Handler) ListStoresByStaffId(c *gin.Context) {
	user_id, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to check user: %w", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to getLimitOffset: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	stores, err := h.s.StaffStoresMapping.ListByStaffID(c.Request.Context(), user_id, lim, off)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to ListByStaffID: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores})
}
