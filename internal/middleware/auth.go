package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MiddlewareService) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, role, err := m.t.ParseClaims(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - " + err.Error()})
			return
		}

		c.Set("id", fmt.Sprintf("%d", id))
		c.Set("role", role)
		c.Next()
	}
}
