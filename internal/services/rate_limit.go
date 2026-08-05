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

func (s *RateLimitService) IsRateLimited(ip string, endpoint models.RateLimitEndpoints) (bool, *models.RateLimitModel, error) {
	limit, err := s.repo.GetRateLimit(ip, endpoint)
	if err != nil && err != sql.ErrNoRows {
		return false, limit, err
	}

	// No record found so not rate limited
	if limit == nil {
		return false, nil, nil
	}

	if limit.RequestCount <= 3 {
		return false, limit, nil
	}

	if t, _ := getTimeRemaining(endpoint, limit); t == 0 {
		return false, limit, nil
	}

	return true, limit, nil
}

func (s *RateLimitService) NewRateLimit(ip string, endpoint models.RateLimitEndpoints) *models.RateLimitModel {
	return &models.RateLimitModel{
		IpAddress:    ip,
		Endpoint:     endpoint,
		RequestCount: 1,
		LastRequest:  time.Now(),
	}
}

func (s *RateLimitService) IncrementRequestCount(limit *models.RateLimitModel, ip string, endpoint models.RateLimitEndpoints) error {
	if limit == nil {
		limit = s.NewRateLimit(ip, endpoint)
	} else {
		limit.RequestCount++
		limit.LastRequest = time.Now()
	}

	return s.repo.UpsertRateLimit(limit)
}

func (s *RateLimitService) ApplyRateLimit(endpoint models.RateLimitEndpoints, limit *models.RateLimitModel) (time.Duration, error) {
	// if there is still time remaining, it means the user is already rate limited
	timeRemaining, waitTime := getTimeRemaining(endpoint, limit)
	if timeRemaining > 0 {
		return timeRemaining, nil
	}

	// if time is twice the wait time that means this is a good behaving user so reset their rate limit
	if time.Since(limit.LastRequest) >= 2*waitTime {
		limit.RequestCount = 1
		limit.LastRequest = time.Now()
		return 0, s.repo.UpsertRateLimit(limit)
	}

	// not much time has passed since last request so increase request count
	// cap at 5 to prevent infinite growth but allow for significant penalties
	if limit.RequestCount >= 5 {
		if limit.Endpoint == models.LimitSendCode && limit.RequestCount <= 10 {
			limit.RequestCount += 2
		}
	} else {
		limit.RequestCount++
	}
	limit.LastRequest = time.Now()
	return 0, s.repo.UpsertRateLimit(limit)
}

func (s *RateLimitService) Clean() {
	for endpoint, multiplier := range models.LimitMultipliers {
		s.repo.CleanExpired(endpoint, multiplier)
	}
}

// Calculate the wait time using the attempts count and the multiplier
func getTimeRemaining(endpoint models.RateLimitEndpoints, limit *models.RateLimitModel) (time.Duration, time.Duration) {
	multiplier := models.LimitMultipliers[endpoint]
	waitTime := time.Duration(math.Pow(2, float64(limit.RequestCount-3))) * time.Minute * time.Duration(multiplier)
	if time.Since(limit.LastRequest) < waitTime {
		timeRemaining := waitTime - time.Since(limit.LastRequest)
		return timeRemaining, waitTime
	}
	return 0, waitTime
}
