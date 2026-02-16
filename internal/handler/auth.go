package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) Login(c *gin.Context) {
	var auth model.Auth

	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&auth); err != nil {
		h.handleFieldError(c, "Login", err)
		return
	}

	u, err := h.s.Authorize.GetByLogin(ctx, auth.Login)
	if err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.handleError(c, "Login", err, http.StatusInternalServerError)
		return
	}

	access_token, refresh_token, err := h.s.Authorize.Authenticate(auth.Password, u)
	if err != nil {
		h.handleError(c, "Login", err, http.StatusUnauthorized)
		return
	}

	err = h.s.Authorize.UpsertRefreshToken(ctx, u.ID, u.Type, refresh_token)
	if err != nil {
		h.handleError(c, "Login - upsert refresh token", err, http.StatusInternalServerError)
		return
	}

	userInfo, err := h.s.Authorize.GetLoginUserInfo(ctx, u)
	if err != nil {
		if errors.Is(err, errs.ErrorUserDeactivated) {
			c.JSON(http.StatusGone, gin.H{"error": "user is deactivated"})
			return
		}
		h.handleError(c, "Login - get user info", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         *userInfo,
	})

}

func (h *Handler) RegisterSystem(c *gin.Context) {
	var user model.CreateSystemUser

	if err := c.ShouldBindJSON(&user); err != nil {
		h.handleFieldError(c, "RegisterSystem", err)
		return
	}

	existingUser, err := h.s.Authorize.GetByLogin(c, user.Login)
	if err != nil && err != errs.ErrorNotFound {
		h.handleError(c, "RegisterSystem", err, http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists with that login"})
		return
	}

	if err := h.s.Authorize.CreateSystem(c, &user); err != nil {
		h.handleError(c, "RegisterStaff", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created system user"})
}
func (h *Handler) RegisterStaff(c *gin.Context) {
	var user model.CreateStaffUser

	if err := c.ShouldBindJSON(&user); err != nil {
		h.handleFieldError(c, "RegisterStaff", err)
		return
	}

	existingUser, err := h.s.Authorize.GetByLogin(c, user.Login)
	if err != nil && err != errs.ErrorNotFound {
		h.handleError(c, "RegisterStaff", err, http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists with that login"})
		return
	}

	if err := h.s.Authorize.CreateStaff(c, &user); err != nil {
		h.handleError(c, "RegisterStaff", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created staff user"})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleFieldError(c, "Refresh", err)
		return
	}

	access, refresh, err := h.s.Authorize.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, errs.ErrorUserDeactivated) {
			c.JSON(http.StatusGone, gin.H{"error": "user is deactivated", "status": "deactivated"})
			return
		}
		h.handleError(c, "Refresh", err, http.StatusUnauthorized)
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
		h.handleFieldError(c, "Logout", err)
		return
	}

	if err := h.s.Authorize.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no found refresh"})
			return
		}
		h.handleError(c, "Logout", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
