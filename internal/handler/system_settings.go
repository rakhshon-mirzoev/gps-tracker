package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) CreateSystemSetting(c *gin.Context) {
	var ss *model.SystemSettings

	user_id, err := checkSystemUser(c)
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

	if err := h.s.SystemSettings.Create(c.Request.Context(), ss); err != nil {
		h.handleError(c, "CreateSystemSetting", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	settings := constants.Settings

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "CreateSystemSetting",
		TargetTable: &settings,
		IPAddress:   &ip,
	}); err != nil {
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

	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "DeleteSystemSetting: check system", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "DeleteSystemSetting: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := h.s.SystemSettings.Delete(c.Request.Context(), parsedID); err != nil {
		h.handleError(c, "DeleteSystemSetting", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	settings := constants.Settings

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "DeleteSystemSetting",
		TargetTable: &settings,
		TargetID:    &parsedID,
		IPAddress:   &ip,
	}); err != nil {
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

	user_id, err := checkSystemUser(c)
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

	if err := h.s.SystemSettings.Update(c.Request.Context(), parsedID, dto); err != nil {
		h.handleError(c, "UpdateSystemSetting", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	settings := constants.Settings

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "UpdateSystemSetting",
		TargetTable: &settings,
		TargetID:    &parsedID,
		IPAddress:   &ip,
	}); err != nil {
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
