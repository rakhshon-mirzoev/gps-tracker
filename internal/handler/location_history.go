package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) Location(c *gin.Context) {
	var loc *model.LocationHistory
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "Location: check staff user", err, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&loc); err != nil {
		h.handleFieldError(c, "Location: bind", err)
		return
	}

	loc.StaffID = user_id

	err = h.s.Location.Create(c.Request.Context(), loc, user_id)
	if err != nil {
		h.handleError(c, "Location: create", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "location created: " + loc.ID.String()})
}

func (h *Handler) AllStaffOnlineOffline(c *gin.Context) {
	id, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "AllStaffOnlineOffline: check system user", err, http.StatusInternalServerError)
		return
	}

	beginStr := c.Query("begin")
	endStr := c.Query("end")

	begin, err := time.Parse("2006-01-02", beginStr)
	if err != nil {
		h.handleFieldError(c, "AllStaffOnlineOffline: parse begin", err)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		h.handleFieldError(c, "AllStaffOnlineOffline: parse end", err)
		return
	}
	end = end.AddDate(0, 0, 1)

	data, err := h.s.Location.GetAllStaffOnlineOffline(c.Request.Context(), id, begin, end)
	if err != nil {
		h.handleError(c, "AllStaffOnlineOffline: get durations", err, http.StatusInternalServerError)
		return
	}

	result := make([]gin.H, 0, len(data))
	for staffID, durations := range data {
		result = append(result, gin.H{
			"staff_id": staffID,
			"online":   durations["online"].Minutes(),
			"offline":  durations["offline"].Minutes(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
