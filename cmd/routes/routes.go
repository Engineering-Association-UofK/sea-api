package routes

import (
	"sea-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

var (
	UserHandler  *handlers.UserHandler
	EventHandler *handlers.EventHandler
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(handlers.LoggingMiddleware())

	r.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"working": "yes"}) })

	api := r.Group("/api")

	user := api.Group("/user")
	user.GET("/all", handlers.RequireRole("ROLE_SUPER_ADMIN"), UserHandler.GetAll)

	event := api.Group("/event")
	event.Use(handlers.AuthMiddleware())
	event.GET("", EventHandler.GetAllEvents)
	event.GET("/:id", EventHandler.GetEventByID)
	event.POST("", EventHandler.CreateEvent)
	event.PUT("", EventHandler.UpdateEvent)
	event.DELETE("/:id", EventHandler.DeleteEvent)

	return r
}
