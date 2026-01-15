package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MiddlewareService) CorsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Headers, Access-Control-Allow-Origin, Access-Control-Allow-Methods, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
