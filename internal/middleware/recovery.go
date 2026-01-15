package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *MiddlewareService) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %v", r)
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"error": "internal server error"},
				)
			}
		}()

		c.Next()
	}
}
