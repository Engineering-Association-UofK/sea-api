package routes

import (
	"sea-api/internal/handlers"
	"sea-api/internal/response"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	UserHandler        *handlers.UserHandler
	EventHandler       *handlers.EventHandler
	MailHandler        *handlers.MailHandler
	CertificateHandler *handlers.CertificateHandler
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		response.InternalServerError(c)
		c.Abort()
	}), gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(handlers.ErrorHandlerMiddleware())

	r.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"working": "yes"}) })

	api := r.Group("/api")

	// ==== Mail ====
	mail := api.Group("/mail")
	mail.POST("", MailHandler.SendMail)

	// ==== Users ====
	user := api.Group("/user")
	user.GET("/all", handlers.RequireRole("ROLE_SUPER_ADMIN"), UserHandler.GetAll)

	// ==== Events ====
	event := api.Group("/event")
	event.Use(handlers.AuthMiddleware())
	event.GET("", EventHandler.GetAllEvents)
	event.GET("/:id", EventHandler.GetEventByID)
	event.POST("", EventHandler.CreateEvent)
	event.PUT("", EventHandler.UpdateEvent)
	event.DELETE("/:id", EventHandler.DeleteEvent)

	// ==== Certificates ====
	certificate := api.Group("/certificate")
	certificate.GET("/verify/:hash", CertificateHandler.VerifyCertificate)
	certificate.GET("/download/:id", CertificateHandler.GetCertificates)
	certificate.Use(handlers.AuthMiddleware())
	certificate.POST("/workshop", CertificateHandler.CreateWorkshopCertificate)
	certificate.GET("/generate-all", CertificateHandler.MakeCertificatesForEvent)
	certificate.GET("/send-all", CertificateHandler.SendCertificatesEmailsForEvent)

	return r
}
