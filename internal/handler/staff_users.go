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
		logger.Log.Errorf("Profile: failed to get lim and off: %w", err)
		return
	}

	bio, err := h.s.StaffUser.GetProfile(c.Request.Context(), user_id, lim, off)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Errorf("Profile: failed to get profile: %w", err)
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
		logger.Log.Errorf("AttachDevice: failed bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.s.StaffUser.Update(c.Request.Context(), user_id, updateStaff, uuid.Nil)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("AttachDevice: failed update: %w", err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) ListStaffUsers(c *gin.Context) {
	staff, err := h.s.StaffUser.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("ListStaffUsers: failed list: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"staff_users": staff})
}

func (h *Handler) CreateStaffUser(c *gin.Context) {
	var dto model.CreateStaffUser

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("CreateStaffUser: failed bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.StaffUser.Create(c.Request.Context(), dto, user_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("CreateStaffUser: failed create: %w", err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) UpdateStaffUser(c *gin.Context) {
	var dto model.UpdateStaffUser

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("UpdateStaffUser: failed to convert id to uuid: %w", err)
		return
	}

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("UpdateStaffUser: failed bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.StaffUser.Update(c.Request.Context(), uid, dto, user_id); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("UpdateStaffUser: failed update: %w", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeactivateStaffUser(c *gin.Context) {
	var deactivate model.UpdateStaffUser

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("DeactivateStaffUser: failed to convert id to uuid: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&deactivate); err != nil {
		logger.Log.Errorf("DeactivateStaffUser: failed bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.s.StaffUser.Update(c.Request.Context(), uid, deactivate, user_id); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("DeactivateStaffUser: failed update: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deactivation complete"})
}

func (h *Handler) AssignSchedule(c *gin.Context) {
	var dto *model.Schedule

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
