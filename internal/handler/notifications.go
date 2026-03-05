package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) ListNotifications(c *gin.Context) {
	userID, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListNotifications: get limit/offset", err)
		return
	}

	notifs, total, err := h.s.Notification.ListByUser(c.Request.Context(), userID, lim, off)
	if err != nil {
		h.handleError(c, "ListNotifications: list", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(notifs, lim, off, total))
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	userID, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "MarkNotificationRead: check staff user", err, http.StatusForbidden)
		return
	}

	nid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "MarkNotificationRead: parse id", err)
		return
	}

	if err := h.s.Notification.MarkRead(c.Request.Context(), nid, userID); err != nil {
		h.handleError(c, "MarkNotificationRead: mark read", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateNotification(c *gin.Context) {
	_, _, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "CreateNotification: check system user", err, http.StatusForbidden)
		return
	}

	var dto model.CreateNotificationRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "CreateNotification: bind json", err)
		return
	}

	if dto.UserID == nil && dto.GroupID == nil {
		h.handleFieldError(c, "CreateNotification", errs.ErrorUserIDorGroupID)
		return
	}

	if dto.UserID != nil && dto.GroupID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "must be 1 of uuids"})
		return
	}

	ctx := c.Request.Context()

	if dto.UserID != nil {
		err = h.s.Notification.CreateAndPushToUser(
			ctx,
			*dto.UserID,
			dto.Title,
			dto.Body,
			dto.Type,
			dto.Meta,
			nil,
		)
	} else if dto.GroupID != nil {
		err = h.s.Notification.CreateAndPushToGroup(
			ctx,
			*dto.GroupID,
			dto.Title,
			dto.Body,
			dto.Type,
			dto.Meta,
		)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "one of uid must be"})
		return
	}

	if err != nil {
		h.handleError(c, "CreateNotification: service", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
