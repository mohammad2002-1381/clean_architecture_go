package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-ca/internal/app/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(string(service.UserIDKey), claims.UserID)
		c.Set(string(service.RoleKey), claims.Role)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, service.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, service.RoleKey, claims.Role)
		ctx = context.WithValue(ctx, service.TokenKey, tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(string(service.RoleKey))
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found in token"})
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role format"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRoleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}
