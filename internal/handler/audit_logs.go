package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuditLogs(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "AuditLogs", err)
		return
	}

	audit_logs, err := h.s.Audit.List(c.Request.Context(), user_id)
	if err != nil {
		h.handleError(c, "AuditLogs", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"audit_logs": audit_logs})
}
