package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
}

func (h *Handler) handleError(c *gin.Context, operation string, err error, defaultStatus int) {
	if err == nil {
		return
	}

	status := defaultStatus
	errorMsg := err.Error()

	if errors.Is(err, errs.ErrorNotFound) {
		status = http.StatusNotFound
		errorMsg = "resource not found"
	}

	requestID := c.GetString(constants.RequestIDKey)

	logger.Log.Error(
		fmt.Sprintf("%s: failed operation", operation),
		zap.Error(err),
		zap.String("operation", operation),
		zap.String("request_id", requestID),
		zap.Int("status_code", status),
	)

	c.JSON(status, ErrorResponse{
		Error:      errorMsg,
		StatusCode: status,
		RequestID:  requestID,
	})
}

func (h *Handler) handleFieldError(c *gin.Context, operation string, err error) {
	requestID := c.GetString(constants.RequestIDKey)

	logger.Log.Warn(
		fmt.Sprintf("%s: field validation error", operation),
		zap.Error(err),
		zap.String("operation", operation),
		zap.String("request_id", requestID),
	)

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:      err.Error(),
		StatusCode: http.StatusBadRequest,
		RequestID:  requestID,
	})
}
