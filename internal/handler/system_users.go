package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) ListSystemUsers(c *gin.Context) {
	system_users, err := h.s.SystemUsers.List(c.Request.Context())
	if err != nil {
		h.handleError(c, "ListSystemUsers: list users", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"system_users": system_users})
}

func (h *Handler) DeleteSystemUser(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeleteSystemUser: parse id", err)
		return
	}

	if err := h.s.SystemUsers.Delete(c.Request.Context(), uid); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found user"})
			return
		}
		h.handleError(c, "DeleteSystemUser: delete user", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateSystemUser(c *gin.Context) {
	var user model.CreateSystemUser

	id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "CreateSystemUser: check user", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "CreateSystemUser: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		h.handleFieldError(c, "CreateSystemUser: bind user", err)
		return
	}

	ip := c.ClientIP()

	if err := h.s.SystemUsers.Create(c.Request.Context(), &user, id, ip, role.(string)); err != nil {
		h.handleError(c, "CreateSystemUser: create user", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})

}

func (h *Handler) UpdateSystemUser(c *gin.Context) {
	var f model.UpdateSystemUser

	id, err := checkSystemUser(c)
	if err != nil {
		h.handleFieldError(c, "UpdateSystemUser: check user", err)
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateSystemUser: parse id", err)
		return
	}

	if err := c.ShouldBindJSON(&f); err != nil {
		h.handleFieldError(c, "UpdateSystemUser: bind user", err)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "UpdateSystemUser: get role", err)
		return
	}

	if err := h.s.SystemUsers.Update(c.Request.Context(), &f, uid, id, c.ClientIP(), role.(string)); err != nil {
		h.handleError(c, "UpdateSystemUser: update system user", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}
