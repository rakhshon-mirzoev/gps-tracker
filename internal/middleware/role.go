package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *MiddlewareService) Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "id problem in token"})

			c.Abort()
		}

		role, ok := c.Get("role")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role problem in token"})

			c.Abort()
		}

		switch role {
		case "administrator":
			c.Next()
		case "supervisor":

			c.Next()

		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not admin or supervisor"})

			c.Abort()
		}

	}
}
