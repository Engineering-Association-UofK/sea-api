package app

import (
	"fmt"
	"sea-api/internal/handlers"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/analyticsrepo"
	"sea-api/internal/repositories/certrepo"
	"sea-api/internal/repositories/electionrepo"
	"sea-api/internal/repositories/eventrepo"
	"sea-api/internal/services"
	"sea-api/internal/services/analytics"
	"sea-api/internal/services/auth"
	"sea-api/internal/services/bot"
	"sea-api/internal/services/certservice"
	"sea-api/internal/services/electionservice"
	"sea-api/internal/services/eventservice"
	"sea-api/internal/services/forms"
	"sea-api/internal/services/schedular"
	"sea-api/internal/services/storage"
	"sea-api/internal/services/userservice"
	st "sea-api/internal/storage"

	"go.uber.org/dig"
)

type Handlers struct {
	dig.In

	User         *handlers.UserHandler
	Event        *handlers.EventHandler
	Mail         *handlers.MailHandler
	Auth         *handlers.AuthHandler
	Account      *handlers.AccountHandler
	Gallery      *handlers.GalleryHandler
	Cms          *handlers.CmsHandler
	Form         *handlers.FormHandler
	Notification *handlers.NotificationHandler
	Bot          *handlers.BotHandler
	Cert         *handlers.CertificatesHandler
	Analytics    *handlers.AnalyticsHandler
	Election     *handlers.ElectionHandler
}

func BuildContainer() (*dig.Container, error) {
	c := dig.New()

	// -------------------------------------------------------------------------
	// 1. INFRASTRUCTURE & STORAGE
	// -------------------------------------------------------------------------
	if err := c.Provide(st.NewMySQLConnection); err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------
	// 2. REPOSITORIES
	// -------------------------------------------------------------------------
	repos := []any{
		repositories.NewUserRepository,
		repositories.NewSuspensionsRepo,
		repositories.NewVerificationRepo,
		repositories.NewFileRepository,
		repositories.NewGalleryRepository,
		repositories.NewCmsRepository,
		repositories.NewFormRepository,
		repositories.NewRateLimitRepository,
		repositories.NewNotificationRepository,
		repositories.NewBotRepository,
		repositories.NewFeedbackRepository,
		repositories.NewAuthRepository,
		analyticsrepo.NewAnalyticsRepo,
		electionrepo.NewElectionRepo,
		certrepo.NewCertRepository,
		eventrepo.NewEventRepository,
	}

	for _, repo := range repos {
		if err := c.Provide(repo); err != nil {
			return nil, fmt.Errorf("failed to provide repository: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// 3. SERVICES
	// -------------------------------------------------------------------------
	svcs := []any{
		storage.NewS3Service,
		services.NewGalleryService,
		services.NewRateLimitService,
		services.NewNotificationService,
		services.NewFeedbackService,
		bot.NewBotService,
		eventservice.NewEventService,
		services.NewAccountService,
		userservice.NewUserService,
		services.NewMailService,
		auth.NewAuthService,
		services.NewCmsService,
		forms.NewFormService,
		analytics.NewAnalytics,
		certservice.NewCertService,
		electionservice.NewElectionService,
		schedular.NewSchedularService,
	}

	for _, svc := range svcs {
		if err := c.Provide(svc); err != nil {
			return nil, fmt.Errorf("failed to provide service: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// 4. HANDLERS
	// -------------------------------------------------------------------------
	handlerConstructors := []any{
		handlers.NewUserHandler,
		handlers.NewEventHandler,
		handlers.NewMailHandler,
		handlers.NewAuthHandler,
		handlers.NewAccountHandler,
		handlers.NewGalleryHandler,
		handlers.NewCmsHandler,
		handlers.NewFormHandler,
		handlers.NewNotificationHandler,
		handlers.NewBotHandler,
		handlers.NewCertificatesHandler,
		handlers.NewAnalyticsHandler,
		handlers.NewElectionHandler,
	}

	for _, h := range handlerConstructors {
		if err := c.Provide(h); err != nil {
			return nil, fmt.Errorf("failed to provide handler: %w", err)
		}
	}

	return c, nil
}
