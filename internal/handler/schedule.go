package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

// ListSchedulesAdmin handles GET /admin/staff/schedules for admin panel
func (h *Handler) ListSchedulesAdmin(c *gin.Context) {
	_, _, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "ListSchedulesAdmin: check system user", err, http.StatusForbidden)
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ListSchedulesAdmin: get limit/offset", err)
		return
	}

	staffIDStr := c.Query("staff_id")
	storeIDStr := c.Query("store_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	schedules, total, err := h.s.Schedule.ListAll(c.Request.Context(), staffIDStr, storeIDStr, dateFrom, dateTo, lim, off)
	if err != nil {
		h.handleError(c, "ListSchedulesAdmin: list", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(schedules, lim, off, total))
}

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

	c.JSON(http.StatusOK, model.NewPaginatedResponse(schedules, lim, off, total))

}

func (h *Handler) AssignSchedule(c *gin.Context) {
	var dto *model.CreateSchedule

	adminID, _, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		h.handleFieldError(c, "AssignSchedule: get role", errs.ErrorRoleEmpty)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "AssignSchedule", err)
		return
	}

	ip := c.ClientIP()
	schedule := constants.Schedule

	if err := h.s.Schedule.Create(c.Request.Context(), dto, adminID, ip, schedule, role.(string), "AssignSchedule"); err != nil {
		h.handleError(c, "AssignSchedule: create schedule", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": dto})
}
