package handler

import (
	"fmt"
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/gin-gonic/gin"
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

	switch err {
	case errs.ErrorNotFound:
		status = http.StatusNotFound
		errorMsg = "resource not found"
	case errs.ErrorAlreadyExists:
		status = http.StatusBadRequest
		errorMsg = "entity already exists"
	case errs.ErrorDuplicatePhotoKey:
		status = http.StatusBadRequest
		errorMsg = "duplicate photo key"
	case errs.ErrorEmptyRequestBody:
		status = http.StatusBadRequest
		errorMsg = "request body is empty"
	case errs.ErrorNotFinishedDay:
		status = http.StatusBadRequest
		errorMsg = "day is not finished yet"
	case errs.ErrorNotStartedDay:
		status = http.StatusBadRequest
		errorMsg = "day is not started yet"
	case errs.ErrorUnprocessable:
		status = http.StatusUnprocessableEntity
		errorMsg = "unporcessable"
	case errs.ErrroActiveVisitExists:
		status = http.StatusBadRequest
		errorMsg = "active visit exists"
	case errs.ErrorWrongPasswordFormat:
		status = http.StatusBadRequest
		errorMsg = "invalid password format"
	case errs.ErrorAlreadyReviewed:
		status = http.StatusBadRequest
		errorMsg = "already reviewed visit"
	case errs.ErrorVisitImmutable:
		status = http.StatusConflict
		errorMsg = "visit status is final and cannot be changed"
	case errs.ErrorUserDeactivated:
		status = http.StatusGone
		errorMsg = "user is deactivated"
	case errs.ErrorInvalidFileFormat:
		status = http.StatusBadRequest
		errorMsg = "file is invalid"
	case errs.ErrorUserIDorGroupID:
		status = http.StatusBadRequest
		errorMsg = "user_id or group_id required"
	case errs.ErrorUniqueViolation:
		status = http.StatusBadRequest
		errorMsg = "sent data duplicates in database"
	case errs.ErrorWeakPassword:
		status = http.StatusBadRequest
		errorMsg = "password is weak"
	case errs.ErrorInvalidPhone:
		status = http.StatusBadRequest
		errorMsg = "phone number is invalid"
	case errs.ErrorInvalidINN:
		status = http.StatusBadRequest
		errorMsg = "INN is invalid"
	case errs.ErrorInvalidInput:
		status = http.StatusBadRequest
		errorMsg = "invalid input"
	case errs.ErrorNotCompletedVisit:
		status = http.StatusBadRequest
		errorMsg = "visit is not completed yet"
	case errs.ErrorScheduleIsNotApproved:
		status = http.StatusBadRequest
		errorMsg = "schedule's not approved"
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
