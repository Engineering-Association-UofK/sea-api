package schedular

import (
	"log/slog"
	"sea-api/internal/handlers/middleware"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
	"time"
)

type SchedularService struct {
	UserRepo         *repositories.UserRepository
	VerificationRepo *repositories.VerificationRepo
	SuspensionsRepo  *repositories.SuspensionsRepo
	BotRepo          *repositories.BotRepository
	MailService      *services.MailService
	RateLimitService *services.RateLimitService
}

func NewSchedularService(
	userRepo *repositories.UserRepository,
	verificationRepo *repositories.VerificationRepo,
	suspensionsRepo *repositories.SuspensionsRepo,
	botRepo *repositories.BotRepository,
	mailService *services.MailService,
	rateLimitService *services.RateLimitService,
) *SchedularService {
	return &SchedularService{
		UserRepo:         userRepo,
		VerificationRepo: verificationRepo,
		SuspensionsRepo:  suspensionsRepo,
		BotRepo:          botRepo,
		MailService:      mailService,
		RateLimitService: rateLimitService,
	}
}

func (s *SchedularService) Run() {
	go s.cleanUpAuthCodes(24 * time.Hour)
	go s.cleanUpSuspensions(2 * time.Hour)
	go s.cleanUpRateLimits(time.Hour)
	go s.cleanUpBotSessions(5 * time.Hour)
}

// For removing codes sent in emails for authentication
func (s *SchedularService) cleanUpAuthCodes(duration time.Duration) {
	cleanCodeTicker := time.NewTicker(duration)
	slog.Debug("codes cleanup service started", "duration", duration)

	s.VerificationRepo.Clean()
	for range cleanCodeTicker.C {
		s.VerificationRepo.Clean()
	}
}

func (s *SchedularService) cleanUpSuspensions(duration time.Duration) {
	cleanSuspensionsTicker := time.NewTicker(duration)
	slog.Debug("suspensions cleanup service started", "duration", duration)

	for range cleanSuspensionsTicker.C {
		ids, err := s.SuspensionsRepo.CleanExpired()
		if err != nil {
			slog.Error("error cleaning expired suspensions", "error", err)
		} else {
			for _, id := range ids {
				err := s.UserRepo.RemoveSuspensionState(id)
				if err != nil {
					slog.Error("error deleting suspension", "error", err, "user_id", id)
				}
			}
		}
	}
}

func (s *SchedularService) cleanUpRateLimits(duration time.Duration) {
	ticker := time.NewTicker(duration)
	slog.Debug("rate limits cleanup service started", "duration", duration)

	for range ticker.C {
		s.RateLimitService.Clean()

		middleware.LimiterMu.Lock()
		for ip, limiter := range middleware.Limiters {
			if limiter.Tokens() == float64(limiter.Burst())*0.85 {
				delete(middleware.Limiters, ip)
			}
		}
		middleware.LimiterMu.Unlock()
	}
}

func (s *SchedularService) cleanUpBotSessions(duration time.Duration) {
	ticker := time.NewTicker(duration)
	slog.Debug("bot sessions cleanup service started", "duration", duration)

	for range ticker.C {
		if err := s.BotRepo.ClearSessions(); err != nil {
			slog.Error("error cleaning expired bot sessions", "error", err)
		}
	}
}
