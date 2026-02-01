package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
)

func (h *Handler) ListStoresByStaffIdAdmin(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		logger.Log.Errorf("ListStoresByStaffIdAdmin: failed to get id: %w", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	h.listStoresByStaffId(c, uid)
}

func (h *Handler) ListStoresByStaffIdStaff(c *gin.Context) {
	uid, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("ListStoresByStaffIdStaff: failed to checkStaffUser: %w", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	h.listStoresByStaffId(c, uid)
}

func (h *Handler) listStoresByStaffId(c *gin.Context, id uuid.UUID) {

	lim, off, err := getLimitOffset(c)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to getLimitOffset: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	stores, err := h.s.StaffStoresMapping.ListByStaffID(c.Request.Context(), id, lim, off)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to ListByStaffID: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores})
}
