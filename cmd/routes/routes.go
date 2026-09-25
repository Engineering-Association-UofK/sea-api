package routes

import (
	"net/http"
	"time"

	"sea-api/cmd/app"
	"sea-api/internal/config"
	"sea-api/internal/handlers/middleware"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"sea-api/internal/services/userservice"

	_ "sea-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

var (
	basicLimit  = middleware.RateLimiter(rate.Every(time.Second), 5)
	midLimit    = middleware.RateLimiter(rate.Every(30*time.Second), 3)
	highLimit   = middleware.RateLimiter(rate.Every(time.Minute), 3)
	strictLimit = middleware.RateLimiter(rate.Every(time.Minute), 1)
)

func SetupRouter(
	u *userservice.UserService,
	rateLimitService *services.RateLimitService,
	h app.Handlers,
) *gin.Engine {
	r := gin.New()
	{ // ==== Config ====
		r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			response.BaseErrorResponse(500, "Internal Server Error", c)
			c.Abort()
		}), gin.Logger())
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}))
		r.Use(middleware.ErrorHandlerMiddleware())
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		r.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": 200}) })

		r.StaticFS("/docs", http.Dir("./docs"))
		r.GET("/docs-ui", func(c *gin.Context) { c.File("./scalar.html") })

		r.GET("/favicon.ico", func(c *gin.Context) { c.File(config.App.ResourcesDir + "/favicon.ico") })
	}
	apiV1 := r.Group("/api/v1")
	apiV1.Use(basicLimit)

	{ // ==== CERTIFICATES
		cert := apiV1.Group("/cert")
		cert.GET("/verify/:hash", h.Cert.VerifyCertificate)
	}

	{ // ==== AUTHENTICATION
		auth := apiV1.Group("/auth")
		auth.POST("/send-verification-code", middleware.StatefulRateLimiter(models.LimitSendCode, rateLimitService), h.Auth.SendVerificationCode)
		auth.POST("/verify", h.Auth.Verify)
		auth.POST("/login", highLimit, h.Auth.Login)
		auth.POST("/forgot-password", highLimit, h.Auth.ForgotPassword)

		auth.POST("/register/check", highLimit, h.Auth.CheckState)
		auth.POST("/register/step", highLimit, h.Auth.DoRegistrationStep)

		auth.POST("/check-username", h.Account.CheckUsernameAvailability)
	}

	{ // ==== EVENTS
		event := apiV1.Group("/event")
		event.GET("", h.Event.GetEventViewList)
		event.GET("/:id", h.Event.GetEventView)
	}

	{ // ==== OPEN
		cms := apiV1.Group("/cms")
		cms.GET("/blogs/:slug", h.Cms.GetViewPostBySlug)
		cms.GET("/blogs", h.Cms.GetViewPostsList)
		cms.GET("/team", h.Cms.GetViewTeamMembers)

		open := apiV1.Group("/open")
		open.POST("/bot", h.Bot.GetNodeView)
	}

	{ // ==== ACCOUNT
		account := apiV1.Group("/account")
		account.Use(middleware.AuthMiddleware(u))

		{ // ==== PROFILE
			account.GET("/summary", h.Account.GetProfileSummary)
			account.GET("", h.Account.GetProfile)
			account.PUT("", h.Account.UpdateProfile)
			account.GET("/certificates", h.Account.GetCertificates)
			// FIXME
			// account.GET("/download/:hash", midLimit, CertificateHandler.GetCertificate)
			account.PUT("/picture", h.Account.UpdatePicture)
			account.PUT("/password", h.Account.UpdatePassword)
			account.PUT("/email", middleware.StatefulRateLimiter(models.LimitUpdateEmail, rateLimitService), h.Account.UpdateEmail)
			account.PUT("/username", middleware.StatefulRateLimiter(models.LimitUpdateUsername, rateLimitService), h.Account.UpdateUsername)
		}

		{ // ==== EVENTS
			event := account.Group("/event")

			// FIXME
			// event.GET("/all-status", h.Event.GetApplicationStatus)
			// event.GET("/status/:id", h.Event.GetOneApplicationStatus)

			event.POST("/:id", h.Event.ApplyForEvent)

			// FIXME
			// event.POST("/cancel/:id", h.Event.CancelApplicationForEvent)
		}

		{ // ==== FORMS
			form := account.Group("/form")
			form.GET("", h.Form.GetAllForms)
			form.GET("/:id", h.Form.GetEntireForUserForm)
		}

		{ // ==== Notifications
			notification := account.Group("/notifications")
			notification.POST("/demo", h.Notification.CreateDemoNotifications)
			notification.GET("", h.Notification.GetNotifications)
			notification.POST("/:id", h.Notification.MarkAsRead)
			notification.POST("", h.Notification.MarkAllAsRead)
			notification.DELETE("/:id", h.Notification.DeleteNotification)
		}
	}

	{ // ###### Administration Endpoints ######
		admin := apiV1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(u), middleware.RequireRole(models.RoleSystemAdmin))

		{ // ==== Analysis
			analysis := admin.Group("/analysis")
			analysis.GET("", midLimit, h.Analytics.GetGeneralAnalytics)
		}

		{ // ==== USERS
			user := admin.Group("/user")
			user.Use(middleware.RequireAnyRole(models.RoleSystemUserMgr, models.RoleSystemSuperAdmin))
			user.GET("/:id", h.User.GetByID)
			user.GET("/all", h.User.GetAll)
			user.POST("/temp-users", h.User.GetAllTempUsers)
			user.GET("/username/:username", h.User.GetByUsername)
			user.POST("/passcode/create/:id", h.User.CreateTempUser)
			user.GET("/passcode/:id", h.User.GetTempUserPasscode)
			user.PUT("", h.User.Update)
			user.POST("/suspend", h.User.Suspend)
			user.POST("/assign-passcodes", h.User.AssignPasscodes)
			user.POST("/import-users-with-emails", h.User.UpdateUsersImport)
			user.POST("/import-users/:id", h.User.ImportUsers)
		}

		{ // ==== ADMIN
			admin.Use(middleware.RequireAnyRole(models.RoleSystemAdminManager, models.RoleSystemSuperAdmin))
			admin.GET("", h.User.GetAdmins)
			admin.POST("/:id", h.User.MakeAdmin)
			admin.PUT("", h.User.UpdateAdmin)
			admin.DELETE("/:id", h.User.DeleteAdmin)
			admin.POST("/add-manager/:id", middleware.RequireRole(models.RoleSystemSuperAdmin), h.User.MakeAdminManager)
			admin.DELETE("/remove-manager/:id", middleware.RequireRole(models.RoleSystemSuperAdmin), h.User.RemoveAdminManager)
		}

		{ // ==== BLOG POSTS
			posts := admin.Group("/blog")
			posts.Use(middleware.RequireAnyRole(models.RoleContentBlogMgr, models.RoleSystemSuperAdmin))
			posts.GET("", h.Cms.GetAllPosts)
			posts.GET("/:id", h.Cms.GetPostById)
			posts.POST("", h.Cms.CreatePost)
			posts.PUT("", h.Cms.UpdatePost)
			posts.DELETE("/:id", h.Cms.DeletePost)
		}

		{ // ==== BOT
			bot := admin.Group("/bot")
			bot.Use(middleware.RequireAnyRole(models.RoleContentEditor, models.RoleSystemSuperAdmin))
			bot.GET("/graph", h.Bot.GetBotGraph)
			bot.PUT("/graph", h.Bot.UpdateBotGraph)
			bot.POST("/reset", h.Bot.ResetDefault)
		}

		{ // ==== GALLERY
			gallery := admin.Group("/gallery")
			gallery.Use(middleware.RequireAnyRole(models.RoleContentEditor, models.RoleSystemSuperAdmin))
			gallery.POST("", h.Gallery.Upload)
			gallery.GET("", h.Gallery.GetAll)
			gallery.GET("/:id", h.Gallery.GetByID)
			gallery.DELETE("", h.Gallery.CleanGallery)
		}

		{ // ==== FORMS
			form := admin.Group("/form")
			form.Use(middleware.RequireAnyRole(models.RoleContentFormMgr, models.RoleSystemSuperAdmin))

			form.GET("", h.Form.GetAllForms)
			form.POST("", h.Form.CreateForm)
			form.PUT("", h.Form.UpdateForm)
			form.DELETE("/:id", h.Form.DeleteForm)

			form.POST("/page", h.Form.CreatePage)
			form.PUT("/page", h.Form.UpdatePage)
			form.DELETE("/page/:id", h.Form.DeletePage)

			form.POST("/question", h.Form.CreateQuestion)
			form.PUT("/question", h.Form.UpdateQuestion)
			form.DELETE("/question/:id", h.Form.DeleteQuestion)

			form.GET("/:id", h.Form.GetEntireForEditForm)

			// FIXME
			// form.POST("/submit", h.Form.SubmitForm)

			form.GET("/analysis/:id", h.Form.GetFormAnalysis)
			form.GET("/detailed-responses/:id", h.Form.GetFormDetailedResponses)

			form.POST("/publish/:id", h.Form.PublishForm)
			form.POST("/unpublish/:id", h.Form.UnpublishForm)

			// form.GET("/user-response/:id", h.Form.GetResponseByID)
			// form.GET("/user-responses/:id", h.Form.GetUserResponsesForForm)

			// form.GET("/responses/:id", h.Form.GetResponsesByFormID)
			// form.PUT("/response-status", h.Form.UpdateResponseStatus)
			// form.DELETE("/response/:id", h.Form.DeleteResponse)

		}

		{ // ==== TEAM MEMBERS
			team := admin.Group("/team")
			team.Use(middleware.RequireAnyRole(models.RoleContentEditor, models.RoleSystemSuperAdmin))
			team.POST("", h.Cms.CreateTeamMember)
			team.GET("", h.Cms.GetAllTeamMembers)
			team.GET("/:id", h.Cms.GetTeamMemberByID)
			team.PUT("", h.Cms.UpdateTeamMember)
			team.DELETE("/:id", h.Cms.DeleteTeamMember)
		}

		{ // ==== EVENTS
			event := admin.Group("/event")
			event.Use(middleware.RequireAnyRole(models.RoleContentEventMgr, models.RoleSystemSuperAdmin))
			event.GET("/:id", h.Event.GetEvent)
			event.GET("", h.Event.GetEventList)
			event.POST("", h.Event.CreateEvent)
			event.PUT("", h.Event.UpdateEvent)
			event.DELETE("/:id", h.Event.DeleteEvent)

			// Coordinators
			event.GET("/:id/coord", h.Event.GetCoordList)
			event.POST("/:id/coord", h.Event.CreateCoords)
			event.PUT("/:id/coord/:coord_id", h.Event.UpdateCoord)
			event.DELETE("/:id/coord/:coord_id", h.Event.DeleteCoord)
			event.DELETE("/:id/coord", h.Event.DeleteCoords)

			// Participation
			event.GET("/:id/application", h.Event.GetApplicationList)
			event.POST("/:id/application/:application_id", h.Event.AcceptApplication)
			event.DELETE("/:id/application/:application_id", h.Event.RejectApplication)

			event.GET("/:id/participant", h.Event.GetParticipantList)
			event.DELETE("/:id/participant/:participation_id", h.Event.RemoveParticipant)
		}

		{ // ==== CERTIFICATES
			certificate := admin.Group("/certificate")
			certificate.Use(middleware.RequireAnyRole(models.RoleCertifier, models.RoleSystemSuperAdmin))

			// New API

			certificate.GET("", basicLimit, h.Cert.GetCertificateList)
			certificate.POST("", midLimit, h.Cert.IssueCertificate)
			certificate.PUT("", midLimit, h.Cert.UpdateCertificate)
			certificate.GET("/:id", midLimit, h.Cert.DownloadCertificate)

			certificate.GET("/template/:id", midLimit, h.Cert.GetTemplate)
			certificate.GET("/template", midLimit, h.Cert.GetTemplatesList)
			certificate.POST("/template", midLimit, h.Cert.CreateTemplate)
			certificate.PUT("/template", midLimit, h.Cert.UpdateTemplate)

			certificate.POST("/test", midLimit, h.Cert.TestGeneration)
		}

		{ // ==== MAIL
			// mail := apiV1.Group("/mail")
			// mail.POST("", MailHandler.SendMail)
		}
	}
	return r
}
