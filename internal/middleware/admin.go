package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *MiddlewareService) Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "id problem in token"})

			c.Abort()
		}

		role, ok := c.Get("role")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role problem in token"})

			c.Abort()
		}

		if role == "administrator" || role == "supervisor" {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not admin"})
		}

	}
}
