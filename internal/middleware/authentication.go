package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m Middleware) Authenticate() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "authorization header is required",
				},
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid authorization header",
				},
			)
			return
		}

		token := parts[1]

		claims, err := m.jwtService.ValidateToken(token)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid or expired token",
				},
			)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
