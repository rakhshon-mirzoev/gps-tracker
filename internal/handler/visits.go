package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) VisitCheckIn(c *gin.Context) {
	var visit model.CreateVisit
	user_id, err := checkStaffUser(c)
	if err != nil {
		h.handleError(c, "VisitCheckIn: check staff user", err, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&visit); err != nil {
		h.handleFieldError(c, "VisitCheckIn", err)
		return
	}

	err = h.s.Visit.Create(c.Request.Context(), visit, user_id)
	if err != nil {
		if errors.Is(err, errs.ErrorAlreadyReviewed) || errors.Is(err, errs.ErrroActiveVisitExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already reviewed or active visit"})
			return
		}
		h.handleError(c, "VisitCheckIn: create visit", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "check-in completed"})

}

func (h *Handler) VisitPhoto(c *gin.Context) {
	var dto model.UpdateVisit

	vID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "VisitPhoto: parse id", err)
		return
	}

	file, err := c.MultipartForm()
	if err != nil {
		h.handleFieldError(c, "VisitPhoto: multipart form", err)
		return
	}

	err = h.s.Visit.VisitPhoto(c.Request.Context(), vID, dto, file)
	if err != nil {
		h.handleError(c, "VisitPhoto: update visit", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update complete"})

}

func (h *Handler) VisitComplete(c *gin.Context) {
	var dto model.UpdateVisit

	staff_id, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "VisitComplete: check staff user", err)
		return
	}

	vID := c.Param("id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleFieldError(c, "VisitComplete", err)
		return
	}

	err = h.s.Visit.VisitComplete(c.Request.Context(), vID, dto, staff_id)
	if err != nil {
		h.handleError(c, "VisitComplete: update visit", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "visit complete"})
}

func (h *Handler) ReviewQueue(c *gin.Context) {
	id, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ReviewQueue: get limit/offset", err)
		return
	}

	status := c.Query("status")

	visits, total, err := h.s.Visit.ListByReviewStatus(c.Request.Context(), status, lim, off, id)
	if err != nil {
		h.handleError(c, "ReviewQueue: list visits", err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, model.NewPaginatedResponse(visits, lim, off, total))
}

func (h *Handler) ReviewVisit(c *gin.Context) {
	var reviewDTO struct {
		Status        string `json:"status" binding:"required"`
		ReviewComment string `json:"review_comment"`
	}

	reviewerID, err := checkSystemUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role empty"})
		return
	}

	visitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleFieldError(c, "ReviewVisit: parse id", err)
		return
	}

	if err := c.ShouldBindJSON(&reviewDTO); err != nil {
		h.handleFieldError(c, "ReviewVisit", err)
		return
	}

	if reviewDTO.Status != "approved" && reviewDTO.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'approved' or 'rejected'"})
		return
	}

	now := time.Now()
	dto := model.UpdateVisit{
		ReviewStatus:  &reviewDTO.Status,
		ReviewComment: &reviewDTO.ReviewComment,
		ReviewedAt:    &now,
	}

	if err := h.s.Visit.UpdateVisit(c.Request.Context(), visitID, dto, reviewerID, role.(string), c.ClientIP(), constants.VisitGroup, "ReviewVisit"); err != nil {
		h.handleError(c, "ReviewVisit: update visit", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review completed", "status": reviewDTO.Status})
}

func (h *Handler) ExportReports(c *gin.Context) {
	var staffID, storeID *uuid.UUID

	if staffIDStr := c.Query("staff_id"); staffIDStr != "" {
		id, err := uuid.Parse(staffIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff_id"})
			return
		}
		staffID = &id
	}

	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store_id"})
			return
		}
		storeID = &id
	}

	var begin, end *time.Time

	if beginStr := c.Query("begin"); beginStr != "" {
		t, err := time.Parse("2006-01-02", beginStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid begin (expected YYYY-MM-DD)"})
			return
		}
		begin = &t
	}

	if endStr := c.Query("end"); endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end (expected YYYY-MM-DD)"})
			return
		}
		t = t.AddDate(0, 0, 1)
		end = &t
	}

	visits, err := h.s.Visit.GetDataForExcel(c.Request.Context(), staffID, storeID, begin, end)
	if err != nil {
		h.handleError(c, "ExportReports: get data", err, http.StatusInternalServerError)
		return
	}

	data, err := h.s.Excel.ExportReports(visits)
	if err != nil {
		h.handleError(c, "ExportReports: export excel", err, http.StatusInternalServerError)
		return
	}

	filename := "visits_report.xlsx"

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Header("Content-Length", strconv.Itoa(len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (h *Handler) ListVisitsByStaffID(c *gin.Context) {
	uid, err := checkStaffUser(c)
	if err != nil {
		h.handleFieldError(c, "ListVisitsByStaffID: check staff user", err)
		return
	}

	lim, off, err := getLimitOffset(c)
	if err != nil {
		h.handleFieldError(c, "ReviewQueue: get limit/offset", err)
		return
	}

	visits, total, _, err := h.s.Visit.ListByStaffID(c.Request.Context(), uid, lim, off)
	if err != nil {
		h.handleError(c, "ListVisitsByStaffID: list visits", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(visits, lim, off, total))
}
