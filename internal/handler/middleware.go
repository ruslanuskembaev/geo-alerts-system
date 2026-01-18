package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет API-key в заголовке
func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" || key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized - valid API key required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
