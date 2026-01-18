package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) Login(c *gin.Context) {
	var auth model.Auth

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

	access_token, refresh_token, err := h.s.Authorize.Authenticate(auth.Password, u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	err = h.s.Authorize.UpsertRefreshToken(ctx, u.ID, u.Type, refresh_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to create refresh_token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access_token, "refresh_token": refresh_token})

}

func (h *Handler) Register(c *gin.Context) {
	var user model.CreateUser

	utype := c.Param("type")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existUser, err := h.s.Authorize.GetByLogin(c.Request.Context(), user.Login)
	if err != nil && err != errs.ErrorNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists user with that login"})
		return
	}

	err = h.s.Authorize.Create(c.Request.Context(), &user, utype)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created user"})

}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.s.Authorize.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.s.Authorize.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if err == errs.ErrorNotFound {
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
