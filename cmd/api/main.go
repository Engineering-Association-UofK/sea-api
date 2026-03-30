package main

import (
	"log/slog"
	"sea-api/cmd/routes"
	"sea-api/internal/config"
	"sea-api/internal/handlers"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
	"sea-api/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	Init()

	r := routes.SetupRouter()
	slog.Info("Starting server on port " + config.App.Port)
	err = r.Run("0.0.0.0:" + config.App.Port)
	if err != nil {
		panic(err)
	}
}

func Init() {
	gin.SetMode(gin.ReleaseMode)
	logger := config.NewMultiHandlerLog()
	slog.SetDefault(logger)

	// Initialize database
	db := storage.NewMySQLConnection()

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	suspensionsRepo := repositories.NewSuspensionsRepo(db)
	eventRepository := repositories.NewEventRepository(db)
	certificateRepository := repositories.NewCertificateRepository(db)
	verificationRepo := repositories.NewVerificationRepo(db)
	fileRepo := repositories.NewFileRepository(db)
	galleryRepository := repositories.NewGalleryRepository(db)
	CmsRepository := repositories.NewCmsRepository(db)
	formRepository := repositories.NewFormRepository(db)

	// Initialize services
	s3StorageService := services.NewS3Service(fileRepo)
	galleryService := services.NewGalleryService(galleryRepository, s3StorageService)
	pdfService := services.NewPDFService(10)

	eventService := services.NewEventService(eventRepository, userRepository)
	accountService := services.NewAccountService(userRepository, s3StorageService)

	userService := services.NewUserService(userRepository, suspensionsRepo, s3StorageService)
	mailService := services.NewMailService(userService)
	authService := services.NewAuthService(userRepository, mailService, verificationRepo)

	CmsService := services.NewCmsService(CmsRepository, userService, galleryService)
	FormService := services.NewFormService(formRepository, galleryService)

	certificateService := services.NewCertificateService(userRepository, eventService, s3StorageService, pdfService, mailService, certificateRepository)
	schedularService := services.NewSchedularService(userRepository, verificationRepo, suspensionsRepo, mailService)
	schedularService.Run()

	// Initialize handlers
	routes.UserHandler = handlers.NewUserHandler(userService)
	routes.EventHandler = handlers.NewEventHandler(eventService)
	routes.MailHandler = handlers.NewMailHandler(mailService)
	routes.CertificateHandler = handlers.NewCertificateHandler(certificateService)
	routes.AuthHandler = handlers.NewAuthHandler(authService)
	routes.AccountHandler = handlers.NewAccountHandler(accountService)
	routes.GalleryHandler = handlers.NewGalleryHandler(galleryService)
	routes.CmsHandler = handlers.NewCmsHandler(CmsService)
	routes.FormHandler = handlers.NewFormHandler(FormService)
}
