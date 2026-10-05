package auth

import (
	"sea-api/internal/config"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/userrepo"
	"sea-api/internal/services"
	"time"
)

type AuthService struct {
	UserRepo         *userrepo.UserRepository
	MailService      *services.MailService
	VerificationRepo *repositories.VerificationRepo
	AuthRepository   *repositories.AuthRepository

	SecretKey  []byte
	Issuer     string
	ExpiryTime time.Duration
}

func NewAuthService(userRepo *userrepo.UserRepository, mailService *services.MailService, verificationRepo *repositories.VerificationRepo, AuthRepository *repositories.AuthRepository) *AuthService {
	return &AuthService{
		UserRepo:         userRepo,
		MailService:      mailService,
		VerificationRepo: verificationRepo,
		AuthRepository:   AuthRepository,
		SecretKey:        []byte(config.App.JwtSecret),
		Issuer:           "SEA - UofK API Server Authorization",
		ExpiryTime:       time.Hour * 168,
	}
}
