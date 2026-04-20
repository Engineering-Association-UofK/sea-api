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
	st "sea-api/internal/services/storage"
	"sea-api/internal/services/user"
	"sea-api/internal/storage"

	"github.com/gin-gonic/gin"
)

// @title						SEA Backend API
// @version					1.0
// @description				This is the backend API for the Steering Engineering Association.
// @contact.name				Technical Office - SEA - UofK
// @contact.email				tech.sea.uofk@gmail.com
// @license.name				MIT
// @license.url				http://opensource.org/licenses/MIT
//
// @host						api-sea-uofk.duckdns.org
// @BasePath					/api/v1
// @schemes					https
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
	S3 := st.NewS3Service(fileRepo)
	galleryService := services.NewGalleryService(galleryRepository, S3)
	rateLimitService := services.NewRateLimitService(rateLimitRepository)
	collaboratorService := services.NewCollaboratorService(collaboratorRepository, S3)
	notificationService := services.NewNotificationService(notificationRepository)

	eventService := services.NewEventService(notificationService, eventRepository, collaboratorRepository, userRepository)
	accountService := services.NewAccountService(userRepository, S3, certificateRepository)

	userService := user.NewUserService(userRepository, suspensionsRepo, S3)
	mailService := services.NewMailService(userService)
	authService := services.NewAuthService(userRepository, mailService, verificationRepo)

	CmsService := services.NewCmsService(CmsRepository, userService, galleryService)
	FormService := forms.NewFormService(formRepository, galleryService)

	certificateService := services.NewCertificateService(userRepository, eventService, S3, pdfService, mailService, collaboratorService, certificateRepository, documentRepository)
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
