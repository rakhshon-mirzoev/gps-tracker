package handler

import (
	"errors"
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateSystemSetting(c *gin.Context) {
	var ss *model.SystemSettings

	user_id, _, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "CreateSystemSetting: check system", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "CreateSystemSetting: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&ss); err != nil {
		h.handleFieldError(c, "CreateSystemSetting", err)
		return
	}

	if err := h.s.SystemSettings.Create(c.Request.Context(), ss, user_id, role.(string), c.ClientIP(), constants.Settings, "CreateSetting"); err != nil {
		h.handleError(c, "CreateSystemSetting", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + ss.ID.String()})

}

func (h *Handler) ListSystemSettings(c *gin.Context) {
	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListSystemSettings: get limit/offset", err)
		return
	}

	ss, total, err := h.s.SystemSettings.List(c.Request.Context(), lim, off)
	if err != nil {
		h.handleError(c, "ListSystemSettings", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(ss, lim, off, total))
}

func (h *Handler) DeleteSystemSetting(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.handleFieldError(c, "DeleteSystemSetting: parse id", err)
		return
	}

	user_id, _, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "DeleteSystemSetting: check system", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "DeleteSystemSetting: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := h.s.SystemSettings.Delete(c.Request.Context(), parsedID, user_id, role.(string), c.ClientIP(), constants.Settings, "DeleteSetting"); err != nil {
		h.handleError(c, "DeleteSystemSetting", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateSystemSetting(c *gin.Context) {
	var dto model.UpdateSystemSettings

	parsedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateSystemSetting: parse id", err)
		return
	}

	user_id, _, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "UpdateSystemSetting: check system", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "UpdateSystemSetting: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateSystemSetting", err)
		return
	}

	if err := h.s.SystemSettings.Update(c.Request.Context(), parsedID, dto, user_id, role.(string), c.ClientIP(), constants.Settings, "UpdateSetting"); err != nil {
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
