package services

import (
	"database/sql"
	"math"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"
)

type RateLimitService struct {
	repo *repositories.RateLimitRepository
}

func NewRateLimitService(repo *repositories.RateLimitRepository) *RateLimitService {
	return &RateLimitService{
		repo: repo,
	}
}

func (s *RateLimitService) IsRateLimited(ip string, endpoint models.RateLimitEndpoints) (bool, time.Duration, error) {
	limit, err := s.repo.GetRateLimit(ip, endpoint)
	if err != nil && err != sql.ErrNoRows {
		return false, 0, err
	}

	// No record found so not rate limited
	if limit == nil {
		return false, 0, s.repo.UpsertRateLimit(&models.RateLimitModel{
			IpAddress:    ip,
			Endpoint:     endpoint,
			RequestCount: 1,
			LastRequest:  time.Now(),
		})
	}

	// Calculate the wait time using the attempts count and the multiplier
	multiplier := models.LimitMultipliers[endpoint]
	waitTime := time.Duration(math.Pow(2, float64(limit.RequestCount-1))) * time.Minute * time.Duration(multiplier)
	if time.Since(limit.LastRequest) < waitTime {
		timeRemaining := waitTime - time.Since(limit.LastRequest)
		return true, timeRemaining, nil
	}

	// if time is twice the wait time that means this is a good behaving user so reset their rate limit
	if time.Since(limit.LastRequest) >= 2*waitTime {
		limit.RequestCount = 1
		limit.LastRequest = time.Now()
		return false, 0, s.repo.UpsertRateLimit(limit)
	}

	// not much time has passed since last request so increase request count
	if limit.RequestCount >= 5 {
		if limit.Endpoint == models.LimitSendCode && limit.RequestCount <= 10 {
			limit.RequestCount += 2
		}
	} else {
		limit.RequestCount++
	}
	limit.LastRequest = time.Now()
	return false, 0, s.repo.UpsertRateLimit(limit)
}

func (s *RateLimitService) Clean() {
	for endpoint, multiplier := range models.LimitMultipliers {
		s.repo.CleanExpired(endpoint, multiplier)
	}
}
