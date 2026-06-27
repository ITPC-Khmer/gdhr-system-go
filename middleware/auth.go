package middleware

import (
	"net/http"
	"strings"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates the Bearer JWT and injects user info into the context.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization token required"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminOnly restricts a route to users with the admin role.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "admin access required"})
			return
		}
		c.Next()
	}
}
