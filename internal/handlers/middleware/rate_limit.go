package middleware

import (
	"fmt"
	"log/slog"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type PenaltyRecord struct {
	Attempts   int
	LastUnlock time.Time
	Multiplier int
}

var (
	Limiters  = make(map[string]*rate.Limiter)
	LimiterMu sync.Mutex
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		LimiterMu.Lock()
		limiter, exists := Limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			Limiters[ip] = limiter
		}
		LimiterMu.Unlock()

		if !limiter.Allow() {
			timeRemaining := limiter.Reserve().DelayFrom(time.Now())
			response.TooManyRequestsErrorResponse(timeRemaining, getMessage(timeRemaining), c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func StatefulRateLimiter(endpoint models.RateLimitEndpoints, service *services.RateLimitService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		isRateLimited, timeRemaining, err := service.IsRateLimited(ip, endpoint)
		if err != nil {
			slog.Error("error checking rate limit", "error", err)
			ctx.Next()
			return
		}
		if isRateLimited {
			response.TooManyRequestsErrorResponse(timeRemaining, getMessage(timeRemaining), ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func getMessage(timeRemaining time.Duration) string {
	msg := "Too many requests. Try again in %s"
	if timeRemaining < time.Minute {
		msg = fmt.Sprintf(msg, timeRemaining.Round(time.Second)) + " seconds."
	} else if timeRemaining < time.Hour {
		msg = fmt.Sprintf(msg, timeRemaining.Round(time.Minute)) + " minutes."
	} else {
		msg = fmt.Sprintf(msg, timeRemaining.Round(time.Hour)) + " hours."
	}
	return msg
}
