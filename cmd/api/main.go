package main

import (
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
	err = r.Run(":" + config.App.Port)
	if err != nil {
		panic(err)
	}
}

func Init() {
	gin.SetMode(gin.ReleaseMode)
	db := storage.NewMySQLConnection()

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	eventRepository := repositories.NewEventRepository(db)
	storeRepository := repositories.NewStoreRepository(db)
	certificateRepository := repositories.NewCertificateRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepository)
	eventService := services.NewEventService(eventRepository, userRepository)
	storageService := services.NewSeaweedService(storeRepository)
	pdfService := services.NewPDFService(10)
	mailService := services.NewMailService(userService)
	certificateService := services.NewCertificateService(userRepository, eventService, storageService, pdfService, mailService, certificateRepository)

	// Initialize handlers
	routes.UserHandler = handlers.NewUserHandler(userService)
	routes.EventHandler = handlers.NewEventHandler(eventService)
	routes.MailHandler = handlers.NewMailHandler(mailService)
	routes.CertificateHandler = handlers.NewCertificateHandler(certificateService)
}
