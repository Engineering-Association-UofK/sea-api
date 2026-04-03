package middleware

import (
	"sea-api/internal/models"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(role models.Role) gin.HandlerFunc {

	return func(c *gin.Context) {

		userData, exists := c.Get("user")

		if !exists {
			c.AbortWithStatus(401)
			return
		}

		claims := userData.(*models.ManagedClaims)

		if !slices.Contains(claims.Roles, role) {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}

func RequireAnyRole(roles ...models.Role) gin.HandlerFunc {

	return func(c *gin.Context) {

		userData, exists := c.Get("user")

		if !exists {
			c.AbortWithStatus(401)
			return
		}

		claims := userData.(*models.ManagedClaims)

		hasRole := false
		for _, role := range roles {
			if slices.Contains(claims.Roles, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}
