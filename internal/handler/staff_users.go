package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
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
		h.handleFieldError(c, "Profile: get limit/offset", err)
		return
	}

	bio, err := h.s.StaffUser.GetProfile(c.Request.Context(), user_id, lim, off)
	if err != nil {
		if errors.Is(err, errs.ErrorUserDeactivated) {
			c.JSON(410, gin.H{"error": "user is deactivated"})
			return
		}
		h.handleError(c, "Profile: get profile", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": bio})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var dto model.UpdateStaffUser

	user_id, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateProfile: bind update staff", err)
		return
	}

	err = h.s.StaffUser.Update(c.Request.Context(), user_id, dto, uuid.Nil, "", "", "", "")
	if err != nil {
		if errors.Is(err, errs.ErrorUserDeactivated) {
			c.JSON(410, gin.H{"error": "user is deactivated"})
			return
		}
		h.handleError(c, "UpdateProfile: update profile", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	userID, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)

	file, err := c.FormFile("avatar")
	if err != nil {
		h.handleFieldError(c, "UpdateAvatar: get file", err)
		return
	}

	err = h.s.StaffUser.UpdateAvatar(c.Request.Context(), userID, file)
	if err != nil {
		h.handleError(c, "UpdateAvatar", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "avatar updated"})
}

func (h *Handler) AttachDevice(c *gin.Context) {
	var updateStaff model.UpdateStaffUser
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "AttachDevice: check staff user", err)
		return
	}

	if err := c.ShouldBindJSON(&updateStaff); err != nil {
		h.handleFieldError(c, "AttachDevice", err)
		return
	}

	err = h.s.StaffUser.Update(c.Request.Context(), user_id, updateStaff, uuid.Nil, "", "", "", "")
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		h.handleError(c, "AttachDevice: update", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ListStaffUsers(c *gin.Context) {
	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListStaffUsers: get limit/offset", err)
		return
	}

	staff, total, err := h.s.StaffUser.List(c.Request.Context(), lim, off)
	if err != nil {
		h.handleError(c, "ListStaffUsers", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(staff, lim, off, total))
}

func (h *Handler) CreateStaffUser(c *gin.Context) {
	var dto model.CreateStaffUser

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "failed to check system user"})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "CreateStaffUser: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "CreateStaffUser", err)
		return
	}

	if err := h.s.StaffUser.Create(c.Request.Context(), dto, user_id, c.ClientIP(), role.(string), constants.UserTypeStaff, "CreateStaffUser"); err != nil {
		h.handleError(c, "CreateStaffUser", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) UpdateStaffUser(c *gin.Context) {
	var dto model.UpdateStaffUser

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateStaffUser: parse id", err)
		return
	}

	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "failed to check system user"})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "UpdateStaffUser: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateStaffUser", err)
		return
	}

	if err := h.s.StaffUser.Update(c.Request.Context(), uid, dto, user_id, c.ClientIP(), role.(string), constants.UserTypeStaff, "UpdateStaffUser"); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
			return
		}
		h.handleError(c, "UpdateStaffUser", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeactivateStaffUser(c *gin.Context) {
	var deactivate model.DeactivateStaffUser

	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "DeactivateStaffUser", err)
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeactivateStaffUser", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "DeactivateStaffUser: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&deactivate); err != nil {
		h.handleError(c, "DeactivateStaffUser", err, http.StatusBadRequest)
		return
	}

	deactivate.DeactivatedByID = &user_id

	if err := h.s.StaffUser.Deactivate(c.Request.Context(), uid, deactivate, user_id, c.ClientIP(), role.(string), constants.UserTypeStaff, "Deactivate"); err != nil {
		h.handleError(c, "DeactivateStaffUser", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deactivation complete"})
}

func (h *Handler) ListStaffUsersLocationByGroupID(c *gin.Context) {
	id, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ListStaffUsersLocationByGroupID: check system user", err, http.StatusBadRequest)
		return
	}

	staff_users, err := h.s.StaffUser.ListByGroupID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, "ListStaffUsersLocationByGroupID: list", err, http.StatusBadRequest)
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

	begin, err := time.Parse("2006-01-02", beginStr)
	if err != nil {
		h.handleFieldError(c, "ListStaffUserLocationByDate: parse begin (expected YYYY-MM-DD)", err)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		h.handleFieldError(c, "ListStaffUserLocationByDate: parse end (expected YYYY-MM-DD)", err)
		return
	}

	end = end.Add(24*time.Hour - time.Nanosecond)

	staff_id, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		h.handleFieldError(c, "ListStaffUserLocationByDate: parse staff id", err)
		return
	}

	staff_user, err := h.s.StaffUser.GetLocationByID(c.Request.Context(), staff_id, begin, end)
	if err != nil {
		h.handleError(c, "ListStaffUserLocationByDate: list locations", err, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"staff_user": staff_user})
}

func (h *Handler) GetStaffUserByID(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "GetStaffUserByID: parse id", err)
		return
	}

	staff_user, err := h.s.StaffUser.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no user found"})
			return
		}
		h.handleError(c, "GetStaffUserByID", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"staff_user": staff_user})
}
