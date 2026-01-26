package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) Profile(c *gin.Context) {
	user_id, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("failed to get lim and off: %w", err)
		return
	}

	bio, err := h.s.StaffUser.GetProfile(c.Request.Context(), user_id, lim, off)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Errorf("failed to get profile: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": bio})
}

func (h *Handler) AttachDevice(c *gin.Context) {
	var updateStaff model.UpdateStaffUser
	user_id, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updateStaff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.s.StaffUser.Update(c.Request.Context(), user_id, updateStaff)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("error in staff update: %w", err)
		return
	}
	c.Status(http.StatusOK)
}

func checkStaffUser(c *gin.Context) (uuid.UUID, error) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	role, _ := c.Get("role")
	if role != "staff" {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	return user_id.(uuid.UUID), nil
}
