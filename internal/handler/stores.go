package handler

import (
	"fmt"
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/IkbolNazarov/gps-tracker/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetStoreByID(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "GetStoreByID: parse id", err)
		return
	}

	store, err := h.s.Store.GetByIDPhoto(c.Request.Context(), uid)
	if err != nil {
		h.handleError(c, "GetStoreByID", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"store": store})
}

func (h *Handler) CreateStore(c *gin.Context) {
	var dto model.CreateStore
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "CreateStore: check staff", err)
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

func (h *Handler) UploadStorePhoto(c *gin.Context) {
	_, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "UploadStorePhoto: check staff", err)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, err := c.FormFile("photo")
	if err != nil {
		h.handleFieldError(c, "UploadStorePhoto: get file", err)
		return
	}

	if err := utils.ValidateFileSize(file); err != nil {
		h.handleFieldError(c, "UploadStorePhoto: validate size", err)
		return
	}

	ext := utils.GetFileExtension(file.Filename)
	objectName := fmt.Sprintf("stores/tmp/%s%s", uuid.New().String(), ext)

	key, err := h.minio.UploadFile(c.Request.Context(), file, objectName)
	if err != nil {
		h.handleError(c, "UploadStorePhoto: upload", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_key": key,
	})
}

func (h *Handler) ListStores(c *gin.Context) {
	isActive := c.Query("is_active")
	status := c.Query("status")

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListStores: get lim off", err)
		return
	}

	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	stores, total, err := h.s.Store.List(c.Request.Context(), isActive, status, sortBy, sortOrder, lim, off)
	if err != nil {
		h.handleError(c, "ListStores: list stores", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(stores, lim, off, total))
}

func (h *Handler) ListStoresByGroupID(c *gin.Context) {
	_, _, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ListStoresByGroupID: check system user", err, http.StatusBadRequest)
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		h.handleFieldError(c, "ListStoresByGroupID: parse group_id", err)
		return
	}

	stores, err := h.s.Store.ListByGroupID(c.Request.Context(), groupID)
	if err != nil {
		h.handleError(c, "ListStoresByGroupID: list stores", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores})
}

func (h *Handler) AssignStoreGroup(c *gin.Context) {
	userID, role, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "AssignStoreGroup: check system user", err, http.StatusBadRequest)
		return
	}

	storeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "AssignStoreGroup: parse store id", err)
		return
	}

	var body struct {
		GroupID uuid.UUID `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.handleFieldError(c, "AssignStoreGroup", err)
		return
	}

	dto := model.UpdateStore{GroupID: &body.GroupID}
	if err := h.s.Store.Update(c.Request.Context(), dto, storeID, userID, c.ClientIP(), constants.Stores, role, "AssignStoreGroup"); err != nil {
		h.handleError(c, "AssignStoreGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "store group assigned"})
}

func (h *Handler) UnassignStoreGroup(c *gin.Context) {
	userID, role, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "UnassignStoreGroup: check system user", err, http.StatusBadRequest)
		return
	}

	storeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UnassignStoreGroup: parse store id", err)
		return
	}

	nilGroupID := uuid.Nil
	dto := model.UpdateStore{GroupID: &nilGroupID}
	if err := h.s.Store.Update(c.Request.Context(), dto, storeID, userID, c.ClientIP(), constants.Stores, role, "UnassignStoreGroup"); err != nil {
		h.handleError(c, "UnassignStoreGroup", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "store group removed"})
}

func (h *Handler) ApproveStore(c *gin.Context) {
	var s model.UpdateStore

	uID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "ApproveStore: parse id", err)
		return
	}

	user_id, _, err := checkSystemUser(c)
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

	user_id, role, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ApproveStore: check system", err, http.StatusBadRequest)
		return
	}

	err = h.s.Store.Delete(c.Request.Context(), uID, user_id, c.ClientIP(), role, constants.Stores, "DeleteStore")
	if err != nil {
		h.handleError(c, "ApproveStore: update store", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
