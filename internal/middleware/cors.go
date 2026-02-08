package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *MiddlewareService) CorsMiddleWare() gin.HandlerFunc {

	return cors.Default()
}
