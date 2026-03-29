package handler

import (
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListStoresByStaffIdAdmin(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		h.handleFieldError(c, "ListStoresByStaffIdAdmin: parse staff_id", err)
		return
	}
	h.listStoresByStaffId(c, uid)
}

func (h *Handler) ListStoresByStaffIdStaff(c *gin.Context) {
	uid, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "ListStoresByStaffIdStaff: check staff user", err, http.StatusForbidden)
		return
	}
	h.listStoresByStaffId(c, uid)
}

func (h *Handler) listStoresByStaffId(c *gin.Context, id uuid.UUID) {
	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListStores: get limit/offset", err)
		return
	}

	stores, total, err := h.s.StaffStoresMapping.ListByStaffID(c.Request.Context(), id, lim, off)
	if err != nil {
		h.handleError(c, "ListStores: list by staff id", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(stores, lim, off, total))
}

func (h *Handler) MapStoreStaff(c *gin.Context) {
	var ssm model.StaffStoresMapping
	if err := c.ShouldBindJSON(&ssm); err != nil {
		h.handleFieldError(c, "MapStoreStaff: bind body", err)
		return
	}

	staff_id, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		h.handleFieldError(c, "MapStoreStaff: parse id", err)
		return
	}

	ssm.StaffID = staff_id

	if err := h.s.StaffStoresMapping.Create(c.Request.Context(), &ssm); err != nil {
		h.handleError(c, "MapStoreStaff", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) UnMapStoreStaff(c *gin.Context) {
	var ssm model.StaffStoresMapping
	if err := c.ShouldBindJSON(&ssm); err != nil {
		h.handleFieldError(c, "UnMapStoreStaff: bind body", err)
		return
	}

	staff_id, err := uuid.Parse(c.Param("staff_id"))
	if err != nil {
		h.handleFieldError(c, "UnMapStoreStaff: parse id", err)
		return
	}

	ssm.StaffID = staff_id

	if err := h.s.StaffStoresMapping.Delete(c.Request.Context(), ssm.StoreID, ssm.StaffID); err != nil {
		h.handleError(c, "UnMapStoreStaff", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
