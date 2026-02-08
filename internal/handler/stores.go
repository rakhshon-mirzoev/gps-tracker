package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) CreateStore(c *gin.Context) {
	var dto *model.CreateStore
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "CreateStore: check staff", err, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "CreateStore", err)
		return
	}

	err = h.s.Store.Create(c.Request.Context(), dto, user_id)
	if err != nil {
		h.handleError(c, "CreateStore: create store", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "create complete"})
}

func (h *Handler) ListStores(c *gin.Context) {
	var f model.StoreListFilter

	if err := c.ShouldBindJSON(&f); err != nil {
		h.handleFieldError(c, "ListStores", err)
		return
	}
	stores, total, err := h.s.Store.List(c.Request.Context(), f)
	if err != nil {
		h.handleError(c, "ListStores: list stores", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores, "total": total})
}

func (h *Handler) ApproveStore(c *gin.Context) {
	var s model.UpdateStore

	uID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "ApproveStore: parse id", err)
		return
	}

	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ApproveStore: check system", err, http.StatusBadRequest)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "ApproveStore: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&s); err != nil {
		h.handleFieldError(c, "ApproveStore", err)
		return
	}

	err = h.s.Store.Update(c.Request.Context(), s, uID, user_id, c.ClientIP(), role.(string), constants.Stores, "ApproveStore")
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no store found"})
			return
		}
		h.handleError(c, "ApproveStore: update store", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteStore(c *gin.Context) {
	uID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "ApproveStore: parse id", err)
		return
	}

	user_id, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ApproveStore: check system", err, http.StatusBadRequest)
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "ApproveStore: get role", errs.ErrorRoleEmpty)
		return
	}

	err = h.s.Store.Delete(c.Request.Context(), uID, user_id, c.ClientIP(), role.(string), constants.Stores, "DeleteStore")
	if err != nil {
		if errors.Is(err, errs.ErrorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no store found"})
			return
		}
		h.handleError(c, "ApproveStore: update store", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
