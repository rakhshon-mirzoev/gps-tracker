package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) Login(c *gin.Context) {
	var (
		auth model.Auth
	)

	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.s.SystemUsers.GetByLogin(ctx, auth.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.s.Authorize.Authenticate(ctx, auth.Login, auth.Password, u.PasswordHash, u.Role, u.ID)
	if err != nil {
		switch err {
		case errors.New("not found"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func (h *Handler) Registrate(c *gin.Context) {
	var (
		user *model.CreateSystemUser
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existUser, err := h.s.SystemUsers.GetByLogin(c.Request.Context(), user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists user with that login"})
		return
	}

	err = h.s.SystemUsers.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created user"})

}
