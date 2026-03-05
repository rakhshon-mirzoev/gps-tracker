package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

func (h *Handler) Location(c *gin.Context) {
	var loc []model.LocationHistory
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "Location: check staff user", err, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&loc); err != nil {
		h.handleFieldError(c, "Location: bind", err)
		return
	}

	if len(loc) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty locations array"})
		return
	}

	for i := range loc {
		loc[i].StaffID = user_id
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.s.Location.CreateBatch(ctx, loc)

	if err != nil {
		h.handleError(c, "Location: create location batch", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "locations created",
		"count":   len(loc),
	})
}

func (h *Handler) AllStaffOnlineOffline(c *gin.Context) {
	id, _, err := checkSystemUser(c)
	if err != nil {
		h.handleError(c, "AllStaffOnlineOffline: check system user", err, http.StatusInternalServerError)
		return
	}

	beginStr := c.Query("begin")
	endStr := c.Query("end")

	if beginStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "begin and end dates are required (expected YYYY-MM-DD)"})
		return
	}

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
