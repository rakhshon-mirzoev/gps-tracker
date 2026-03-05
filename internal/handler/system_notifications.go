package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

// GetSystemNotifications возвращает список уведомлений для авторизованного admin/supervisor.
//
// GET /api/v1/admin/notifications?limit=20&offset=0
func (h *Handler) GetSystemNotifications(c *gin.Context) {
	recipientID, _, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "GetSystemNotifications: get limit/offset", err)
		return
	}

	notifs, total, err := h.s.SystemNotification.ListByRecipient(c.Request.Context(), recipientID, lim, off)
	if err != nil {
		h.handleError(c, "GetSystemNotifications: list", err, http.StatusInternalServerError)
		return
	}

	// Возвращаем в стандартном paginated формате проекта
	type paginatedResponse struct {
		Data   interface{} `json:"data"`
		Total  int         `json:"total"`
		Limit  int         `json:"limit"`
		Offset int         `json:"offset"`
	}

	c.JSON(http.StatusOK, paginatedResponse{
		Data:   notifs,
		Total:  total,
		Limit:  lim,
		Offset: off,
	})
}

// MarkSystemNotificationRead помечает уведомление прочитанным.
//
// PATCH /api/v1/admin/notifications/:id/read
func (h *Handler) MarkSystemNotificationRead(c *gin.Context) {
	recipientID, _, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	nid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "MarkSystemNotificationRead: parse id", err)
		return
	}

	if err := h.s.SystemNotification.MarkRead(c.Request.Context(), nid, recipientID); err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		h.handleError(c, "MarkSystemNotificationRead: mark read", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
