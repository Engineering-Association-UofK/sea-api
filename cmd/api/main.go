package main

import (
	"fmt"
	"log/slog"
	"sea-api/cmd/app"
	"sea-api/cmd/routes"
	"sea-api/internal/config"
	"sea-api/internal/services"
	"sea-api/internal/services/schedular"
	"sea-api/internal/services/userservice"

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

	container, err := app.BuildContainer()
	if err != nil {
		panic(fmt.Sprintf("Failed to build DI container: %v", err))
	}

	// Resolve scheduler and execute background tasks
	err = container.Invoke(func(schedularService *schedular.SchedularService) {
		schedularService.Run()
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to start scheduler: %v", err))
	}

	// Resolve router dependencies and run HTTP server
	err = container.Invoke(func(
		userService *userservice.UserService,
		rateLimitService *services.RateLimitService,
		h app.Handlers,
	) {
		r := routes.SetupRouter(userService, rateLimitService, h)

		slog.Info("Starting server on port " + config.App.Port)
		if err := r.Run("0.0.0.0:" + config.App.Port); err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to invoke application runner: %v", err))
	}
}
