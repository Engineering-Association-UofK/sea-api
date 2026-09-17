package main

import (
	"log/slog"
	"sea-api/cmd/routes"
	"sea-api/internal/config"
	"sea-api/internal/handlers"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/certrepo"
	"sea-api/internal/repositories/eventrepo"
	"sea-api/internal/services"
	"sea-api/internal/services/auth"
	"sea-api/internal/services/bot"
	"sea-api/internal/services/certservice"
	"sea-api/internal/services/eventservice"
	"sea-api/internal/services/forms"
	"sea-api/internal/services/schedular"
	st "sea-api/internal/services/storage"
	"sea-api/internal/services/userservice"
	"sea-api/internal/storage"

	"github.com/gin-gonic/gin"
)

// @title						SEA Backend API
// @version						1.0
// @description					This is the backend API for the Steering Engineering Association.
// @contact.name				Technical Office - SEA - UofK
// @contact.email				tech.sea.uofk@gmail.com
// @license.name				MIT
// @license.url					http://opensource.org/licenses/MIT
//
// @host						api-sea-uofk-dev.duckdns.org
// @BasePath					/api/v1
// @schemes						https
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Type 'Bearer <token>' to authenticate
func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	Go()
}

func Go() {
	gin.SetMode(gin.ReleaseMode)
	logger := config.NewMultiHandlerLog(slog.Level(config.App.LoggingLevel))
	slog.SetDefault(logger)

	// Initialize database
	db := storage.NewMySQLConnection()

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	suspensionsRepo := repositories.NewSuspensionsRepo(db)

	// FIXME: Remove implementation
	// eventRepository := eventrepo.NewEventRepository(db)

	// FIXME: Remove implementation
	// certificateRepository := repositories.NewCertificateRepository(db)
	verificationRepo := repositories.NewVerificationRepo(db)
	fileRepo := repositories.NewFileRepository(db)
	galleryRepository := repositories.NewGalleryRepository(db)
	CmsRepository := repositories.NewCmsRepository(db)
	formRepository := repositories.NewFormRepository(db)
	rateLimitRepository := repositories.NewRateLimitRepository(db)

	// FIXME: Remove implementation
	// documentRepository := repositories.NewDocumentRepository(db)
	notificationRepository := repositories.NewNotificationRepository(db)
	botRepository := repositories.NewBotRepository(db)
	feedbackRepository := repositories.NewFeedbackRepository(db)
	authRepository := repositories.NewAuthRepository(db)
	certRepo := certrepo.NewCertRepository(db)
	eventRepo := eventrepo.NewEventRepository(db)

	// Initialize services

	// FIXME: Remove implementation
	// pdfService := services.NewPDFService(10)
	S3 := st.NewS3Service(fileRepo)
	galleryService := services.NewGalleryService(galleryRepository, S3)
	rateLimitService := services.NewRateLimitService(rateLimitRepository)
	notificationService := services.NewNotificationService(notificationRepository)
	feedbackService := services.NewFeedbackService(feedbackRepository)

	botService := bot.NewBotService(botRepository, feedbackService)
	eventService := eventservice.NewEventService(eventRepo, formRepository, S3, galleryService)
	accountService := services.NewAccountService(userRepository, S3, certRepo)

	userService := userservice.NewUserService(userRepository, suspensionsRepo, S3)
	mailService := services.NewMailService(userService)
	authService := auth.NewAuthService(userRepository, mailService, verificationRepo, authRepository)

	CmsService := services.NewCmsService(CmsRepository, userService, galleryService)
	FormService := forms.NewFormService(formRepository, eventService, galleryService)

	certService := certservice.NewCertService(certRepo, S3, eventService, userService)

	// FIXME: Remove implementation
	// certificateService := cert.NewCertificateService(
	// 	userRepository,
	// 	eventService,
	// 	S3,
	// 	pdfService,
	// 	mailService,
	// 	collaboratorService,
	// 	notificationService,
	// 	certificateRepository,
	// 	documentRepository,
	// )
	schedularService := schedular.NewSchedularService(
		userRepository,
		verificationRepo,
		suspensionsRepo,
		botRepository,
		mailService,
		rateLimitService,
	)
	schedularService.Run()

	// Initialize handlers
	routes.UserHandler = handlers.NewUserHandler(userService)
	routes.EventHandler = handlers.NewEventHandler(eventService)
	routes.MailHandler = handlers.NewMailHandler(mailService)
	// FIXME: Remove implementation
	// routes.CertificateHandler = handlers.NewCertificateHandler(certificateService)
	routes.AuthHandler = handlers.NewAuthHandler(authService)
	routes.AccountHandler = handlers.NewAccountHandler(accountService)
	routes.GalleryHandler = handlers.NewGalleryHandler(galleryService)
	routes.CmsHandler = handlers.NewCmsHandler(CmsService)
	routes.FormHandler = handlers.NewFormHandler(FormService)
	routes.NotificationHandler = handlers.NewNotificationHandler(notificationService)
	routes.BotHandler = handlers.NewBotHandler(botService)
	routes.CertHandler = handlers.NewCertificatesHandler(certService)

	// Initialize routes

	r := routes.SetupRouter(userService, rateLimitService)
	slog.Info("Starting server on port " + config.App.Port)
	err := r.Run("0.0.0.0:" + config.App.Port)
	if err != nil {
		panic(err)
	}
}
