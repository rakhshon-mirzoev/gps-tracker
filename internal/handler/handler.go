package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
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
