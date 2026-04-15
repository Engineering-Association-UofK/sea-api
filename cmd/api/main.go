package main

import (
	"log/slog"
	"sea-api/cmd/routes"
	"sea-api/internal/config"
	"sea-api/internal/handlers"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
	"sea-api/internal/services/forms"
	"sea-api/internal/services/schedular"
	"sea-api/internal/storage"

	"github.com/gin-gonic/gin"
)

// @title						SEA Backend API
// @version					1.0
// @description				This is the backend for the Steering Engineering Association.
// @host						localhost:8000
// @BasePath					/api/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type 'Bearer <token>' to authenticate
func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	Go()
}

func Go() {
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
	collaboratorRepository := repositories.NewCollaboratorRepo(db)
	rateLimitRepository := repositories.NewRateLimitRepository(db)
	documentRepository := repositories.NewDocumentRepository(db)
	notificationRepository := repositories.NewNotificationRepository(db)

	// Initialize services
	pdfService := services.NewPDFService(10)
	s3StorageService := services.NewS3Service(fileRepo)
	galleryService := services.NewGalleryService(galleryRepository, s3StorageService)
	rateLimitService := services.NewRateLimitService(rateLimitRepository)
	collaboratorService := services.NewCollaboratorService(collaboratorRepository, s3StorageService)
	notificationService := services.NewNotificationService(notificationRepository)

	eventService := services.NewEventService(notificationService, eventRepository, userRepository)
	accountService := services.NewAccountService(userRepository, s3StorageService, certificateRepository)

	userService := services.NewUserService(userRepository, suspensionsRepo, s3StorageService)
	mailService := services.NewMailService(userService)
	authService := services.NewAuthService(userRepository, mailService, verificationRepo)

	CmsService := services.NewCmsService(CmsRepository, userService, galleryService)
	FormService := forms.NewFormService(formRepository, galleryService)

	certificateService := services.NewCertificateService(userRepository, eventService, s3StorageService, pdfService, mailService, collaboratorService, certificateRepository, documentRepository)
	schedularService := schedular.NewSchedularService(userRepository, verificationRepo, suspensionsRepo, mailService, rateLimitService)
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
	routes.CollaboratorHandler = handlers.NewCollaboratorHandler(collaboratorService)
	routes.NotificationHandler = handlers.NewNotificationHandler(notificationService)

	r := routes.SetupRouter(userService, rateLimitService)
	slog.Info("Starting server on port " + config.App.Port)
	err := r.Run("0.0.0.0:" + config.App.Port)
	if err != nil {
		panic(err)
	}
}
