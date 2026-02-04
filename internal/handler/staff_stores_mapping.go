package handler

import (
	"net/http"

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

	stores, err := h.s.StaffStoresMapping.ListByStaffID(c.Request.Context(), id, lim, off)
	if err != nil {
		h.handleError(c, "ListStores: list by staff id", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"stores": stores})
}
