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

func (h *Handler) CreateStore(c *gin.Context) {
	var dto *model.CreateStore
	user_id, err := checkStaffUser(c)
	if err != nil {
		logger.Log.Errorf("CreateStore: failed to check staff: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		logger.Log.Errorf("CreateStore: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.s.Store.Create(c.Request.Context(), dto, user_id)
	if err != nil {
		logger.Log.Errorf("CreateStore: failed to create store: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "create complete"})
}

func (h *Handler) ListStores(c *gin.Context) {
	var f model.StoreListFilter

	if err := c.ShouldBindJSON(&f); err != nil {
		logger.Log.Errorf("ListStores: failed to bind: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	stores, total, err := h.s.Store.List(c.Request.Context(), f)
	if err != nil {
		logger.Log.Errorf("ListStores: failed to list stores: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores, "total": total})
}

func (h *Handler) ApproveStore(c *gin.Context) {
	var s model.UpdateStore

	uID, err := uuid.FromBytes([]byte(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("ApproveStore: failed to convert id to uuid: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("ApproveStore: failed to bind: %w", err)
		return
	}

	err = h.s.Store.Update(c.Request.Context(), s, uID)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no store found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("ApproveStore: failed update: %w", err)
		return
	}

	c.Status(http.StatusOK)
}
