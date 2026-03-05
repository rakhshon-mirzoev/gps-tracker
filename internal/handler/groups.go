package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) CreateGroup(c *gin.Context) {
	user_id, _, err := checkSystemUser(c)
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

	if err := h.s.Group.Create(c.Request.Context(), group, user_id, c.ClientIP(), constants.Group, role.(string), "CreateGroup"); err != nil {
		h.handleError(c, "CreateGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + group.ID.String()})
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	user_id, _, err := checkSystemUser(c)
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

	if err := h.s.Group.Delete(c.Request.Context(), uid, user_id, c.ClientIP(), constants.Group, role.(string), "DeleteGroup"); err != nil {
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
		h.handleError(c, "GetByIdGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": g})
}

func (h *Handler) ListGroups(c *gin.Context) {
	userID, role, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ListGroups: check system user", err, http.StatusBadRequest)
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListGroups: get limit/offset", err)
		return
	}

	var supervisorID *uuid.UUID

	if role == "supervisor" {
		supervisorID = &userID
	} else {
		supervisorID = nil
	}

	groups, total, err := h.s.Group.List(c.Request.Context(), supervisorID, lim, off)
	if err != nil {
		h.handleError(c, "ListGroups: list", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(groups, lim, off, total))
}
func (h *Handler) UpdateGroup(c *gin.Context) {
	user_id, _, err := checkSystemUser(c)
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

	if err := h.s.Group.Update(c.Request.Context(), dto, uid, user_id, c.ClientIP(), constants.Group, role.(string), "UpdateGroup"); err != nil {
		h.handleError(c, "UpdateGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update"})
}
