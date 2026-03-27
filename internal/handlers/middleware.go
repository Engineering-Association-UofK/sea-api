package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()
		slog.Debug("Logging middleware", "Method", c.Request.Method, "Path", c.Request.URL.Path, "Time took", time.Since(start))
	}
}

// Checks the validation of the JWT token made by the Spring backend
func (u *UserHandler) AuthMiddleware() gin.HandlerFunc {
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

		claims.Roles, err = u.service.GetRolesByUserID(claims.UserID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err

		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c)
			return
		}

		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case errs.BadRequest:
				response.NewBaseError(http.StatusBadRequest, appErr.Message, c)

			case errs.NotFound:
				response.NewBaseError(http.StatusNotFound, appErr.Message, c)

			case errs.Unauthorized:
				response.NewBaseError(http.StatusUnauthorized, appErr.Message, c)

			case errs.Forbidden:
				response.NewBaseError(http.StatusForbidden, appErr.Message, c)

			case errs.Conflict:
				response.NewBaseError(http.StatusConflict, appErr.Message, c)

			case errs.MultiBadRequest:
				response.NewErrorResponse(http.StatusBadRequest, appErr.Message, c, appErr.Fields)

			default:
				response.InternalServerError(c)
			}
			return
		}

		response.InternalServerError(c)
	}
}
