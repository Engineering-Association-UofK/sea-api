package handlers

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(role string) gin.HandlerFunc {

	return func(c *gin.Context) {

		userData, exists := c.Get("user")

		if !exists {
			c.AbortWithStatus(401)
			return
		}

		claims := userData.(*models.UserClaims)

		if !slices.Contains(claims.Roles, role) {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()
		slog.Debug("Logging middleware", "Method", c.Request.Method, "Path", c.Request.URL.Path, "Time took", time.Since(start))
	}
}

// Checks the validation of the JWT token made by the Spring backend
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secretBytes, err := base64.RawStdEncoding.DecodeString(config.App.JwtSecret)
		if err != nil {
			slog.Error("Failed to decode JWT secret from Base64", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// Get token and claims
		claims := &models.UserClaims{}
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

		c.Set("user", claims)
		c.Next()
	}
}
