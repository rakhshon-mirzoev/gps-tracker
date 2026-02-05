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

func (h *Handler) CreateGroup(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "CreateGroup: CheckSystemUser", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "CreateGroup: GetRole", errs.ErrorRoleEmpty)
		return
	}

	var group *model.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		h.handleFieldError(c, "CreateGroup", err)
		return
	}

	if err := h.s.Group.Create(c.Request.Context(), group); err != nil {
		h.handleError(c, "CreateGroup", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	groups := constants.Group

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "CreateGroup",
		TargetTable: &groups,
		IPAddress:   &ip,
	}); err != nil {
		h.handleError(c, "CreateGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + group.ID.String()})
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "DeleteGroup: CheckSystemUser", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "DeleteGroup: GetRole", errs.ErrorRoleEmpty)
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeleteGroup: parse id", err)
		return
	}

	if err := h.s.Group.Delete(c.Request.Context(), uid); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		h.handleError(c, "DeleteGroup", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	groups := constants.Group

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "DeleteGroup",
		TargetTable: &groups,
		TargetID:    &uid,
		IPAddress:   &ip,
	}); err != nil {
		h.handleError(c, "DeleteGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted group"})
}

func (h *Handler) GetByIdGroup(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "GetByIdGroup: parse id", err)
		return
	}

	g, err := h.s.Group.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		h.handleError(c, "GetByIdGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": g})
}

func (h *Handler) ListGroups(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ListGroups: check system user", err, http.StatusBadRequest)
		return
	}
	groups, err := h.s.Group.List(c.Request.Context(), &user_id)
	if err != nil {
		h.handleError(c, "ListGroups: list", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (h *Handler) UpdateGroup(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "UpdateGroup: CheckSystemUser", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "UpdateGroup: GetRole", errs.ErrorRoleEmpty)
		return
	}

	var dto model.Group

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateGroup: parse id", err)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateGroup", err)
		return
	}

	if err := h.s.Group.Update(c.Request.Context(), dto, uid); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		h.handleError(c, "UpdateGroup", err, http.StatusInternalServerError)
		return
	}

	ip := c.ClientIP()
	groups := constants.Group

	if err := h.s.Audit.Create(c.Request.Context(), &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role.(string),
		Action:      "UpdateGroup",
		TargetTable: &groups,
		TargetID:    &uid,
		IPAddress:   &ip,
	}); err != nil {
		h.handleError(c, "UpdateGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update"})
}
