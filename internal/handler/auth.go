package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
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

	u, err := h.s.Authorize.GetByLogin(ctx, auth.Login)
	if err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.s.Authorize.Authenticate(auth.Password, u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})

}

func (h *Handler) Registrate(c *gin.Context) {
	var (
		user *model.CreateUser
	)

	utype := c.Param("type")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existUser, err := h.s.Authorize.GetByLogin(c.Request.Context(), user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists user with that login"})
		return
	}

	err = h.s.Authorize.Create(c.Request.Context(), user, utype)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created user"})

}
