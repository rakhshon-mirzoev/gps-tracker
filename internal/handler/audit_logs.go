package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) AuditLogs(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "AuditLogs", err)
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "AuditLogs: get limit/offset", err)
		return
	}

	action := c.Query("action")
	targetTable := c.Query("target_table")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	audit_logs, total, err := h.s.Audit.List(c.Request.Context(), user_id, action, targetTable, dateFrom, dateTo, lim, off)
	if err != nil {
		h.handleError(c, "AuditLogs", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(audit_logs, lim, off, total))
}
