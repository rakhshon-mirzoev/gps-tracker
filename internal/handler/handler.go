package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Handler struct {
	s     *service.Service
	t     token.Token
	minio minioservice.Minio
}

func New(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}

func getLimitOffset(c *gin.Context) (int, int, error) {
	lim, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return 0, 0, err
	}

	off, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return 0, 0, err
	}
	return lim, off, nil
}

func checkStaffUser(c *gin.Context) (uuid.UUID, error) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	role, _ := c.Get("role")
	if role != "staff" {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	return user_id.(uuid.UUID), nil
}

func checkSystemUser(c *gin.Context) (uuid.UUID, error) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	role, _ := c.Get("role")
	if role != "admin" && role != "supervisor" {
		return uuid.Nil, errs.ErrorPermissionDenied
	}

	return user_id.(uuid.UUID), nil
}
