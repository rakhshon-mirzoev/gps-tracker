package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

func (h *Handler) VisitCheckIn(c *gin.Context) {
	var visit = &model.CreateVisit{}
	user_id, err := checkStaffUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("VisitCheckIn: failed to check staff user: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("VisitCheckIn: failed to check bind visit: %w", err)
		return
	}

	err = h.s.Visit.Create(c.Request.Context(), visit, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("VisitCheckIn: failed to create visit: %w", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "check-in completed"})

}

func (h *Handler) VisitPhoto(c *gin.Context) {
	var dto model.UpdateVisit

	vID := c.Param("id")

	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении файлов: " + err.Error()})
		logger.Log.Errorf("VisitPhoto: multipart form: %w", err)
		return
	}

	err = h.s.Visit.VisitPhoto(c.Request.Context(), vID, dto, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("VisitPhoto: failed to update visit: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update complete"})

}

func (h *Handler) VisitComplete(c *gin.Context) {
	var dto model.UpdateVisit

	vID := c.Param("id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("VisitComplete: failed to bind dto: %w", err)
		return
	}

	err := h.s.Visit.VisitComplete(c.Request.Context(), vID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("VisitComplete: failed to update visit: %w", err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		logger.Log.Errorf("ReviewQueue: failed to get limit/offset: %w", err)
		return
	}

	visits, total, err := h.s.Visit.ListByReviewStatus(c.Request.Context(), "pending", lim, off, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		logger.Log.Errorf("ReviewQueue: failed to list visits: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"visits": visits, "total": total})
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

	visitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visit id"})
		logger.Log.Errorf("ReviewVisit: failed to parse id: %w", err)
		return
	}

	if err := c.ShouldBindJSON(&reviewDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Errorf("ReviewVisit: failed to bind dto: %w", err)
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
		ReviewerID:    &reviewerID,
		ReviewedAt:    &now,
	}

	if err := h.s.Visit.UpdateVisit(c.Request.Context(), visitID, dto); err != nil {
		if err == errs.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Errorf("ReviewVisit: failed to update visit: %w", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review completed", "status": reviewDTO.Status})
}
