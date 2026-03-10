package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BaseError struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBaseError(status int, msg string, c *gin.Context) {
	c.JSON(status, BaseError{
		Status:    status,
		Message:   msg,
		Timestamp: time.Now(),
	})
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, BaseError{
		Status:    500,
		Message:   "Internal Server Error",
		Timestamp: time.Now(),
	})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, BaseError{
		Status:    404,
		Message:   "Not Found",
		Timestamp: time.Now(),
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, BaseError{
		Status:    401,
		Message:   "Unauthorized",
		Timestamp: time.Now(),
	})
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, BaseError{
		Status:    400,
		Message:   "Bad Request",
		Timestamp: time.Now(),
	})
}
