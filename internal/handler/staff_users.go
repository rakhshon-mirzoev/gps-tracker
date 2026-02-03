package handler

import (
	"errors"
	"net/http"
	"time"

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
		c.JSON(http.StatusForbidden, gin.H{"error": "failed to check system user"})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("CreateStaffUser: failed to bind staff user: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind staff user"})
		return
	}

	if err := h.s.StaffUser.Create(c.Request.Context(), dto, user_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff user"})
		logger.Log.Errorf("CreateStaffUser: failed create: %w", err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) UpdateStaffUser(c *gin.Context) {
	var dto model.UpdateStaffUser

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse uuid"})
		logger.Log.Errorf("UpdateStaffUser: failed to convert id to uuid: %w", err)
		return
	}

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "failed to check system user"})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("UpdateStaffUser: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind update staff user"})
		return
	}

	if err := h.s.StaffUser.Update(c.Request.Context(), uid, dto, user_id); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update staff user"})
		logger.Log.Errorf("UpdateStaffUser: failed update: %w", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeactivateStaffUser(c *gin.Context) {
	var deactivate model.UpdateStaffUser

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "failed to check system user"})
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert id to uuid"})
		logger.Log.Errorf("DeactivateStaffUser: failed to convert id to uuid: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&deactivate); err != nil {
		logger.Log.Errorf("DeactivateStaffUser: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	now := time.Now()
	is_active := false
	deactivate.DeactivatedByID = &user_id
	deactivate.DeactivatedAt = &now
	deactivate.IsActive = &is_active

	if err := h.s.StaffUser.Update(c.Request.Context(), uid, deactivate, user_id); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update"})
		logger.Log.Errorf("DeactivateStaffUser: failed update: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deactivation complete"})
}

func (h *Handler) ListStaffUsersLocationByGroupID(c *gin.Context) {
	id, err := checkSystemUser(c)
	if err != nil {
		logger.Log.Errorf("ListStaffUsersLocationByGroupID: failed to check system user: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to check system user"})
		return
	}

	staff_users, err := h.s.StaffUser.ListByGroupID(c.Request.Context(), id)
	if err != nil {
		logger.Log.Errorf("ListStaffUsersLocationByGroupID: failed to list staff users: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to list staff users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"staff_users": staff_users})
}

func (h *Handler) ListStaffUserLocationByDate(c *gin.Context) {
	beginStr := c.Query("begin")
	endStr := c.Query("end")

	if beginStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "begin and end dates are required"})
		return
	}

	begin, err := time.Parse("02.01.2006", beginStr)
	if err != nil {
		logger.Log.Errorf("ListStaffUserLocationByDate: failed to parse begin time: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid begin date format, use dd.mm.yyyy"})
		return
	}

	end, err := time.Parse("02.01.2006", endStr)
	if err != nil {
		logger.Log.Errorf("ListStaffUserLocationByDate: failed to parse end time: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format, use dd.mm.yyyy"})
		return
	}

	end = end.Add(24*time.Hour - time.Nanosecond)

	staff_id, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		logger.Log.Errorf("ListStaffUserLocationByDate: failed to parse staff id: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse staff id"})
		return
	}

	staff_user, err := h.s.StaffUser.GetLocationByID(c.Request.Context(), staff_id, begin, end)
	if err != nil {
		logger.Log.Errorf("ListStaffUserLocationByDate: failed to list staff user's location history: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to list staff user's location history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"staff_user": staff_user})
}
