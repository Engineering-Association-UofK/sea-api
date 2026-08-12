package middleware

import (
	"fmt"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/services/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(s *user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secretBytes := []byte(config.App.JwtSecret)

		// Get token and claims
		claims := &models.ManagedClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secretBytes, nil
		})

		// Check validity
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		u, err := s.GetByUserID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if !u.Verified {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not verified"})
			c.Abort()
			return
		}

		claims.Roles, err = s.GetRolesByUserID(claims.UserID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
