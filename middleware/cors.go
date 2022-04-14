package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AllowCORSConfig is allowing all incoming requests
func AllowCORSConfig() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = []string{"Location"}
	return cors.New(config)
}
