package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) DeleteVisitPhoto(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeleteVisitPhoto: uuid parse", err)
		return
	}

	if err := h.s.VisitPhotos.Delete(c.Request.Context(), uid); err != nil {
		h.handleError(c, "DeleteVisitPhoto: delete visit photo", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
