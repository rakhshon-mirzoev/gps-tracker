package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) ListShedules(c *gin.Context) {
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "ListShedules: check staff user", err, http.StatusForbidden)
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListShedules: get limit/offset", err)
		return
	}

	schedules, total, err := h.s.Schedule.ListByStaffID(c.Request.Context(), user_id, lim, off)
	if err != nil {
		h.handleError(c, "ListShedules: list by staff id", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules, "total": total})

}

func (h *Handler) AssignSchedule(c *gin.Context) {
	var dto *model.CreateSchedule

	adminID, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "AssignSchedule", err)
		return
	}

	if err := h.s.Schedule.Create(c.Request.Context(), dto, adminID); err != nil {
		h.handleError(c, "AssignSchedule: create schedule", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": dto})
}
