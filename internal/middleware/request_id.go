package middleware

import (
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *MiddlewareService) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(constants.RequestIDKey, reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
