package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
)

type Handler struct {
	S *service.Service
}

func InitHandler(s *service.Service) *Handler {
	return &Handler{S: s}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}
