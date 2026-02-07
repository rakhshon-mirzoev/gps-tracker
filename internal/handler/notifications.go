package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	c.JSON(http.StatusOK, gin.H{"notifications": notifs, "total": total})
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	userID, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	nid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "MarkNotificationRead: parse id", err)
		return
	}

	if err := h.s.Notification.MarkRead(c.Request.Context(), nid, userID); err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		h.handleError(c, "MarkNotificationRead: mark read", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
