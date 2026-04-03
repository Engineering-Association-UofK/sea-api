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

type ErrorResponse struct {
	BaseError
	Errors map[string]string `json:"errors"`
}

type TooManyResponse struct {
	BaseError
	TimeRemaining time.Duration `json:"time_remaining"`
}

func BaseErrorResponse(status int, msg string, c *gin.Context) {
	c.JSON(status, BaseError{
		Status:    status,
		Message:   msg,
		Timestamp: time.Now(),
	})
}

func ListErrorResponse(status int, msg string, c *gin.Context, errors map[string]string) {
	c.JSON(status, ErrorResponse{
		BaseError: BaseError{
			Status:    status,
			Message:   msg,
			Timestamp: time.Now(),
		},
		Errors: errors,
	})
}

func TooManyRequestsErrorResponse(timeRemaining time.Duration, msg string, c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, TooManyResponse{
		BaseError: BaseError{
			Status:    http.StatusTooManyRequests,
			Message:   msg,
			Timestamp: time.Now(),
		},
		TimeRemaining: timeRemaining,
	})
}
