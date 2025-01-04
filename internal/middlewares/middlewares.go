package middlewares

import (
	"authentication/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.Valid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized: tokens expired")
			c.Abort()
		}
		c.Next()
	}
}
