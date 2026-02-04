package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
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
		h.handleFieldError(c, "Register", err)
		return
	}

	_, err := h.s.Authorize.GetByLogin(c.Request.Context(), user.Login)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found user with login"})
		}
		h.handleError(c, "Register", err, http.StatusInternalServerError)
		return
	}

	err = h.s.Authorize.Create(c.Request.Context(), &user, utype)
	if err != nil {
		h.handleError(c, "Register", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": "created user"})

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
