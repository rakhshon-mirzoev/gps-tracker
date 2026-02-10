package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListSystemUsers(c *gin.Context) {
	system_users, err := h.s.SystemUsers.List(c.Request.Context())
	if err != nil {
		h.handleError(c, "ListSystemUsers: list users", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"system_users": system_users})
}
