package handler

import (
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	adminID, role, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "AssignSchedule", err)
		return
	}

	if err := h.s.Schedule.Create(c.Request.Context(), dto, adminID, c.ClientIP(), constants.Schedule, role, "AssignSchedule"); err != nil {
		h.handleError(c, "AssignSchedule: create schedule", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": dto})
}

func (h *Handler) AssignScheduleBulk(c *gin.Context) {
	var dto *model.CreateSchedulesBulk

	adminID, role, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "AssignScheduleBulk", err)
		return
	}

	count, err := h.s.Schedule.CreateBulk(c.Request.Context(), dto, adminID, c.ClientIP(), constants.Schedule, role)
	if err != nil {
		h.handleError(c, "AssignScheduleBulk: create schedules", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "count": count})
}

func (h *Handler) DeleteSchedule(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "DeleteSchedule: id parse", err)
		return
	}

	if err := h.s.Schedule.Delete(c.Request.Context(), uid); err != nil {
		h.handleError(c, "DeleteSchedule: delete schedule", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateSchedule(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "UpdateSchedule: id parse", err)
		return
	}

	var dto model.UpdateSchedule

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "UpdateSchedule: bind dto", err)
		return
	}

	if err := h.s.Schedule.Update(c.Request.Context(), uid, dto); err != nil {
		h.handleError(c, "UpdateStaffSchedule: update schedule", err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) InitSchedule(c *gin.Context) {
	var dto model.CreateSchedule

	staff_id, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "InitSchedule: check staff", err)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "InitSchedule: bind request", err)
		return
	}

	dto.StaffID = staff_id

	if err := h.s.Schedule.InitSchedule(c.Request.Context(), &dto); err != nil {
		h.handleError(c, "InitSchedule: error to init", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "unplanned schedule created"})

}
