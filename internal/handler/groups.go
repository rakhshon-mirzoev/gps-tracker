package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) CreateGroup(c *gin.Context) {
	var group *model.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		logger.Log.Errorf("CreateGroup: failed to bind: %w", err)
		return
	}

	if err := h.s.Group.Create(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
		logger.Log.Errorf("CreateGroup: failed to create: %w", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created: " + group.ID.String()})
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	uid, err := uuid.FromBytes([]byte(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("DeleteGroup: failed to convert id to uuid: %w", err)
		return
	}

	if err := h.s.Group.Delete(c.Request.Context(), uid); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("DeleteGroup: failed to delete: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted group"})
}

func (h *Handler) GetByIdGroup(c *gin.Context) {
	uid, err := uuid.FromBytes([]byte(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("GetByID: failed to convert id to uuid: %w", err)
		return
	}

	g, err := h.s.Group.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("GetByID: failed to get: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": g})
}

func (h *Handler) ListGroups(c *gin.Context) {
	user_id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("ListGroups: failed to check user: %w", err)
		return
	}
	groups, err := h.s.Group.List(c.Request.Context(), &user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("ListGroups: failed to list: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (h *Handler) UpdateGroup(c *gin.Context) {
	var dto model.Group

	uid, err := uuid.FromBytes([]byte(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("UpdateGroup: failed to convert id to uuid: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("UpdateGroup: failed to bind: %w", err)
		return
	}

	if err := h.s.Group.Update(c.Request.Context(), dto, uid); err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not found group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("UpdateGroup: failed to update: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update"})
}
