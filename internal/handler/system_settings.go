package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) CreateSystemSetting(c *gin.Context) {
	var ss *model.SystemSettings

	if err := c.ShouldBindJSON(&ss); err != nil {
		h.handleFieldError(c, "CreateSystemSetting", err)
		return
	}

	if err := h.s.SystemSettings.Create(c.Request.Context(), ss); err != nil {
		h.handleError(c, "CreateSystemSetting", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + ss.ID.String()})

}

func (h *Handler) ListSystemSettings(c *gin.Context) {

	ss, err := h.s.SystemSettings.List(c.Request.Context())
	if err != nil {
		h.handleError(c, "ListSystemSettings", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"system_settings": ss})

}

func (h *Handler) DeleteSystemSetting(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.handleFieldError(c, "DeleteSystemSetting: parse id", err)
		return
	}

	if err := h.s.SystemSettings.Delete(c.Request.Context(), parsedID); err != nil {
		h.handleError(c, "DeleteSystemSetting", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateSystemSetting(c *gin.Context) {
	var dto model.UpdateSystemSettings

	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.handleFieldError(c, "UpdateSystemSetting: parse id", err)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateSystemSetting", err)
		return
	}

	if err := h.s.SystemSettings.Update(c.Request.Context(), parsedID, dto); err != nil {
		h.handleError(c, "UpdateSystemSetting", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetSystemSettingByID(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "GetSystemSettingByID: parse id", err)
		return
	}

	setting, err := h.s.SystemSettings.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		h.handleError(c, "GetSystemSettingByID", err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"setting": setting})
}
