package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"sea-api/internal/errs"
	"sea-api/internal/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err

		if errors.Is(err, sql.ErrNoRows) {
			response.BaseErrorResponse(http.StatusNotFound, "Not Found", c)
			return
		}

		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case errs.BadRequest:
				response.BaseErrorResponse(http.StatusBadRequest, appErr.Message, c)

			case errs.NotFound:
				response.BaseErrorResponse(http.StatusNotFound, appErr.Message, c)

			case errs.Unauthorized:
				response.BaseErrorResponse(http.StatusUnauthorized, appErr.Message, c)

			case errs.Forbidden:
				response.BaseErrorResponse(http.StatusForbidden, appErr.Message, c)

			case errs.Conflict:
				response.BaseErrorResponse(http.StatusConflict, appErr.Message, c)

			case errs.MultiBadRequest:
				response.ListErrorResponse(http.StatusBadRequest, appErr.Message, c, appErr.Fields)

			default:
				response.BaseErrorResponse(http.StatusInternalServerError, "Internal Server Error", c)
			}
			return
		}

		response.BaseErrorResponse(http.StatusInternalServerError, "Internal Server Error", c)
	}
}
