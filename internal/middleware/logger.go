package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
)

func (s *MiddlewareService) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		reqID, _ := c.Get(constants.RequestIDKey)

		log.Printf(
			"[REQ] id=%v method=%s path=%s status=%d duration=%s",
			reqID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
