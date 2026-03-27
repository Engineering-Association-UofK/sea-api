package services

import (
	"log/slog"
	"sea-api/internal/repositories"
	"time"
)

type SchedularService struct {
	UserRepo         *repositories.UserRepository
	VerificationRepo *repositories.VerificationRepo
	SuspensionsRepo  *repositories.SuspensionsRepo
	MailService      *MailService
}

func NewSchedularService(userRepo *repositories.UserRepository, verificationRepo *repositories.VerificationRepo, suspensionsRepo *repositories.SuspensionsRepo, mailService *MailService) *SchedularService {
	return &SchedularService{
		UserRepo:         userRepo,
		VerificationRepo: verificationRepo,
		SuspensionsRepo:  suspensionsRepo,
		MailService:      mailService,
	}
}

func (s *SchedularService) Run() {
	go s.cleanUpCodes(24 * time.Hour)
	go s.cleanUpSuspensions(2 * time.Hour)
}

func (s *SchedularService) cleanUpCodes(duration time.Duration) {
	cleanCodeTicker := time.NewTicker(duration)

	s.VerificationRepo.Clean()
	for range cleanCodeTicker.C {
		s.VerificationRepo.Clean()
	}
}

func (s *SchedularService) cleanUpSuspensions(duration time.Duration) {
	cleanSuspensionsTicker := time.NewTicker(duration)

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
