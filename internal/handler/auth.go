package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
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
		logger.Log.Errorf("Login: get by login err: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	access_token, refresh_token, err := h.s.Authorize.Authenticate(auth.Password, u)
	if err != nil {
		logger.Log.Errorf("Login: failed to authenticate: %w", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	err = h.s.Authorize.UpsertRefreshToken(ctx, u.ID, u.Type, refresh_token)
	if err != nil {
		logger.Log.Errorf("Login: failed to upsert: %w", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to create refresh_token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access_token, "refresh_token": refresh_token})

}

func (h *Handler) RegisterSystem(c *gin.Context) {
	h.register(c, constants.UserTypeSystem)
}

func (h *Handler) RegisterStaff(c *gin.Context) {
	h.register(c, constants.UserTypeStaff)
}

func (h *Handler) register(c *gin.Context, utype string) {
	var user model.CreateUser

	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Log.Errorf("Register: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind create user"})
		return
	}

	existUser, err := h.s.Authorize.GetByLogin(c.Request.Context(), user.Login)
	if err != nil && err != errs.ErrorNotFound {
		logger.Log.Errorf("Register: failed to getbylogin: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user by login"})
		return
	}

	if existUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists user with that login"})
		return
	}

	err = h.s.Authorize.Create(c.Request.Context(), &user, utype)
	if err != nil {
		logger.Log.Errorf("Register: failed to create user: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created user"})

}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Refresh: failed to bind ref: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind request"})
		return
	}

	access, refresh, err := h.s.Authorize.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Log.Errorf("Refresh: failed to refresh: %w", err)
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
		logger.Log.Errorf("Logout: failed to bind ref: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind refresh token"})
		return
	}

	if err := h.s.Authorize.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no found refresh"})
			return
		}
		logger.Log.Errorf("Logout: failed to logout: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.Status(http.StatusOK)
}
