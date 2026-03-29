package handler

import (
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) DeleteStaffSchedule(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeleteStaffSChedule: id parse", err)
		return
	}

	if err := h.s.StaffSchedule.Delete(c.Request.Context(), uid); err != nil {
		h.handleError(c, "DeleteStaffSchedule: delete schedule", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateStaffSchedule(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateStaffSchedule: id parse", err)
		return
	}

	var dto model.UpdateStaffSchedule

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateStaffSchedule: bind dto", err)
		return
	}

	if err := h.s.StaffSchedule.Update(c.Request.Context(), uid, dto); err != nil {
		h.handleError(c, "UpdateStaffSchedule: update schedule", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ListStaffSchedule(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "ListStaffSchedule: id parse", err)
		return
	}

	schedules, err := h.s.StaffSchedule.List(c.Request.Context(), uid)
	if err != nil {
		h.handleError(c, "ListStaffSchedule: list schedules", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}
