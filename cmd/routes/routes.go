package routes

import (
	"time"

	"sea-api/internal/handlers"
	"sea-api/internal/handlers/middleware"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"

	_ "sea-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

var (
	UserHandler         *handlers.UserHandler
	EventHandler        *handlers.EventHandler
	MailHandler         *handlers.MailHandler
	CertificateHandler  *handlers.CertificateHandler
	AuthHandler         *handlers.AuthHandler
	AccountHandler      *handlers.AccountHandler
	GalleryHandler      *handlers.GalleryHandler
	CmsHandler          *handlers.CmsHandler
	FormHandler         *handlers.FormHandler
	CollaboratorHandler *handlers.CollaboratorHandler
	NotificationHandler *handlers.NotificationHandler
)

var (
	basicLimit  = middleware.RateLimiter(rate.Every(time.Second), 5)
	midLimit    = middleware.RateLimiter(rate.Every(30*time.Second), 3)
	highLimit   = middleware.RateLimiter(rate.Every(time.Minute), 3)
	strictLimit = middleware.RateLimiter(rate.Every(time.Minute), 1)
)

func SetupRouter(u *services.UserService, rateLimitService *services.RateLimitService) *gin.Engine {
	r := gin.New()
	{ // ==== Config ====
		r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			response.BaseErrorResponse(500, "Internal Server Error", c)
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
		r.Use(middleware.ErrorHandlerMiddleware())
		r.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": 200}) })
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	apiV1 := r.Group("/api/v1")
	apiV1.Use(basicLimit)

	{ // ==== CERTIFICATES
		cert := apiV1.Group("/cert")
		cert.GET("/verify/:hash", CertificateHandler.VerifyCertificate)
		cert.GET("/download/:hash", midLimit, CertificateHandler.GetCertificates)
		cert.GET("/verify-document/:hash", CertificateHandler.VerifyDocument)
	}

	{ // ==== AUTHENTICATION
		auth := apiV1.Group("/auth")
		auth.POST("/send-verification-code", middleware.StatefulRateLimiter(models.LimitSendCode, rateLimitService), AuthHandler.SendVerificationCode)
		auth.POST("/verify", AuthHandler.Verify)
		auth.POST("/login", highLimit, AuthHandler.Login)
		auth.POST("/register", highLimit, AuthHandler.Register)
		auth.POST("/check-username", AccountHandler.CheckUsernameAvailability)
	}

	{ // ==== CMS
		cms := apiV1.Group("/cms")
		cms.GET("/blogs/:slug", CmsHandler.GetViewPostBySlug)
		cms.GET("/blogs", CmsHandler.GetViewPostsList)
		cms.GET("/team", CmsHandler.GetViewTeamMembers)
	}

	{ // ==== ACCOUNT
		account := apiV1.Group("/account")
		account.Use(middleware.AuthMiddleware(u))

		{ // ==== PROFILE
			account.GET("", AccountHandler.GetProfile)
			account.PUT("", AccountHandler.UpdateProfile)
			account.GET("/certificates", AccountHandler.GetCertificates)
			account.PUT("/picture", AccountHandler.UpdatePicture)
			account.PUT("/password", AccountHandler.UpdatePassword)
			account.PUT("/email", middleware.StatefulRateLimiter(models.LimitUpdateEmail, rateLimitService), AccountHandler.UpdateEmail)
			account.PUT("/username", middleware.StatefulRateLimiter(models.LimitUpdateUsername, rateLimitService), AccountHandler.UpdateUsername)
		}

		{ // ==== EVENTS & FORMS
			event := account.Group("/event")
			event.GET("/form/:id", FormHandler.GetEntireForUserForm) // <------------------- New
		}

		{ // ==== Notifications
			notification := account.Group("/notifications")
			notification.POST("/demo", NotificationHandler.CreateDemoNotifications)
			notification.GET("", NotificationHandler.GetNotifications)
			notification.POST("/:id", NotificationHandler.MarkAsRead)
			notification.POST("", NotificationHandler.MarkAllAsRead)
			notification.DELETE("/:id", NotificationHandler.DeleteNotification)
		}
	}

	{ // ###### Administration Endpoints ######
		admin := apiV1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(u), middleware.RequireRole(models.RoleSystemAdmin))

		{ // ==== USERS
			user := admin.Group("/user")
			user.Use(middleware.RequireAnyRole(models.RoleSystemUserMgr, models.RoleSystemSuperAdmin))
			user.GET("/:id", UserHandler.GetByID)
			user.GET("/all", UserHandler.GetAll)
			user.POST("/temp-users", UserHandler.GetAllTempUsers)
			user.GET("/username/:username", UserHandler.GetByUsername)
			user.GET("/passcode/:id", UserHandler.GetTempUserPasscode)
			user.PUT("", UserHandler.Update)
			user.POST("/suspend", UserHandler.Suspend)
			user.POST("/assign-passcodes", UserHandler.AssignPasscodes)
			user.POST("/import-users-with-emails", UserHandler.UpdateUsersImport)
		}

		{ // ==== ADMIN
			admin.Use(middleware.RequireAnyRole(models.RoleSystemAdminManager, models.RoleSystemSuperAdmin))
			admin.GET("", UserHandler.GetAdmins)
			admin.POST("/:id", UserHandler.MakeAdmin)
			admin.PUT("", UserHandler.UpdateAdmin)
			admin.DELETE("/:id", UserHandler.DeleteAdmin)
			admin.POST("/add-manager/:id", middleware.RequireRole(models.RoleSystemSuperAdmin), UserHandler.MakeAdminManager)
			admin.DELETE("/remove-manager/:id", middleware.RequireRole(models.RoleSystemSuperAdmin), UserHandler.RemoveAdminManager)
		}

		{ // ==== BLOG POSTS
			blog := admin.Group("/blog")
			blog.Use(middleware.RequireAnyRole(models.RoleContentBlogMgr, models.RoleSystemSuperAdmin))
			blog.GET("", CmsHandler.GetAllPosts)
			blog.GET("/:id", CmsHandler.GetPostById)
			blog.POST("", CmsHandler.CreatePost)
			blog.PUT("", CmsHandler.UpdatePost)
			blog.DELETE("/:id", CmsHandler.DeletePost)
		}

		{ // ==== GALLERY
			gallery := admin.Group("/gallery")
			gallery.Use(middleware.RequireAnyRole(models.RoleContentEditor, models.RoleSystemSuperAdmin))
			gallery.POST("", GalleryHandler.Upload)
			gallery.GET("", GalleryHandler.GetAll)
			gallery.GET("/:id", GalleryHandler.GetByID)
			gallery.DELETE("", GalleryHandler.CleanGallery)
		}

		// TODO: Add bot commands

		{ // ==== FORMS
			// form := admin.Group("/form")
			form := apiV1.Group("/form")
			// form.Use(middleware.RequireAnyRole(models.RoleContentFormMgr, models.RoleSystemSuperAdmin))

			form.GET("", FormHandler.GetAllForms)
			form.POST("", FormHandler.CreateForm)
			form.PUT("", FormHandler.UpdateForm)
			form.DELETE("/:id", FormHandler.DeleteForm)

			form.POST("/page", FormHandler.CreatePage)
			form.PUT("/page", FormHandler.UpdatePage)
			form.DELETE("/page/:id", FormHandler.DeletePage)

			form.POST("/question", FormHandler.CreateQuestion)
			form.PUT("/question", FormHandler.UpdateQuestion)
			form.DELETE("/question/:id", FormHandler.DeleteQuestion)

			form.GET("/:id", FormHandler.GetEntireForEditForm)

			form.POST("/submit", FormHandler.SubmitForm)

			form.GET("/analysis/:id", FormHandler.GetFormAnalysis)
			form.GET("/detailed-responses/:id", FormHandler.GetFormDetailedResponses)

			// form.GET("/user-response/:id", FormHandler.GetResponseByID)
			// form.GET("/user-responses/:id", FormHandler.GetUserResponsesForForm)

			// form.GET("/responses/:id", FormHandler.GetResponsesByFormID)
			// form.PUT("/response-status", FormHandler.UpdateResponseStatus)
			// form.DELETE("/response/:id", FormHandler.DeleteResponse)

		}

		{ // ==== TEAM MEMBERS
			team := admin.Group("/team")
			team.Use(middleware.RequireAnyRole(models.RoleContentEditor, models.RoleSystemSuperAdmin))
			team.POST("", CmsHandler.CreateTeamMember)
			team.GET("", CmsHandler.GetAllTeamMembers)
			team.GET("/:id", CmsHandler.GetTeamMemberByID)
			team.PUT("", CmsHandler.UpdateTeamMember)
			team.DELETE("/:id", CmsHandler.DeleteTeamMember)
		}

		{ // ==== EVENTS
			event := admin.Group("/event")
			event.Use(middleware.RequireAnyRole(models.RoleContentEventMgr, models.RoleSystemSuperAdmin))
			event.GET("", EventHandler.GetAllEvents)
			event.GET("/:id", EventHandler.GetEventByID)
			event.POST("", EventHandler.CreateEvent)
			event.PUT("", EventHandler.UpdateEvent)
			event.DELETE("/:id", EventHandler.DeleteEvent)
			event.GET("/send-all-emails", strictLimit, CertificateHandler.SendCertificatesEmailsForEvent)
			event.POST("/import-users/:id", EventHandler.ImportUsers)
		}

		{ // ==== Collaborators
			collabs := admin.Group("/collabs")
			collabs.Use(middleware.RequireAnyRole(models.RoleContentEventMgr, models.RoleSystemSuperAdmin))
			collabs.GET("", CollaboratorHandler.GetAll)
			collabs.GET("/:id", CollaboratorHandler.GetByID)
			collabs.POST("", CollaboratorHandler.Create)
			collabs.PUT("", CollaboratorHandler.Update)
			collabs.DELETE("/:id", CollaboratorHandler.Delete)
		}

		{ // ==== CERTIFICATES
			certificate := admin.Group("/certificate")
			certificate.Use(middleware.RequireAnyRole(models.RoleCertifier, models.RoleSystemSuperAdmin))
			certificate.GET("/generate-all-for-event", strictLimit, CertificateHandler.MakeCertificatesForEvent)
			certificate.POST("/sign", midLimit, CertificateHandler.SignPDF)
		}

		{ // ==== MAIL
			mail := apiV1.Group("/mail")
			mail.POST("", MailHandler.SendMail)
		}
	}
	return r
}
