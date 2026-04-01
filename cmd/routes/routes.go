package routes

import (
	"time"

	h "sea-api/internal/handlers"
	m "sea-api/internal/models"
	"sea-api/internal/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	UserHandler        *h.UserHandler
	EventHandler       *h.EventHandler
	MailHandler        *h.MailHandler
	CertificateHandler *h.CertificateHandler
	AuthHandler        *h.AuthHandler
	AccountHandler     *h.AccountHandler
	GalleryHandler     *h.GalleryHandler
	CmsHandler         *h.CmsHandler
	FormHandler        *h.FormHandler
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	{ // ==== Config ====
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
		r.Use(h.ErrorHandlerMiddleware())
		r.GET("/test", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": 200}) })
	}
	apiV1 := r.Group("/api/v1")

	{ // ==== CERTIFICATES
		cert := apiV1.Group("/cert")
		cert.GET("/verify/:hash", CertificateHandler.VerifyCertificate)
		cert.GET("/download/:hash", CertificateHandler.GetCertificates)
	}

	{ // ==== AUTHENTICATION
		auth := apiV1.Group("/auth")
		auth.POST("/send-verification-code", AuthHandler.SendVerificationCode)
		auth.POST("/verify", AuthHandler.Verify)
		auth.POST("/login", AuthHandler.Login)
		auth.POST("/register", AuthHandler.Register)
	}

	{ // ==== ACCOUNT
		account := apiV1.Group("/account")
		account.Use(UserHandler.AuthMiddleware())

		{ // ==== PROFILE
			account.GET("", AccountHandler.GetProfile)
			account.PUT("", AccountHandler.UpdateProfile)
			account.PUT("/picture", AccountHandler.UpdatePicture)
			account.PUT("/password", AccountHandler.UpdatePassword)
			account.PUT("/email", AccountHandler.UpdateEmail)
			account.PUT("/username", AccountHandler.UpdateUsername)
			account.POST("/check-username", AccountHandler.CheckUsernameAvailability)
		}

		{ // ==== EVENTS & FORMS
			event := account.Group("/event")
			event.GET("/form/:id", FormHandler.GetEntireForUserForm) // <------------------- New
		}
	}

	{ // ###### Administration Endpoints ######
		admin := apiV1.Group("/admin")
		admin.Use(UserHandler.AuthMiddleware(), h.RequireRole(m.RoleSystemAdmin))

		{ // ==== USERS
			user := admin.Group("/user")
			user.Use(h.RequireAnyRole(m.RoleSystemUserMgr, m.RoleSystemSuperAdmin))
			user.GET("/:id", UserHandler.GetByID)
			user.POST("/all", UserHandler.GetAll)
			user.POST("/temp-users", UserHandler.GetAllTempUsers)
			user.GET("/username/:username", UserHandler.GetByUsername)
			user.GET("/passcode/:id", UserHandler.GetTempUserPasscode)
			user.PUT("", UserHandler.Update)
			user.POST("/suspend", UserHandler.Suspend)
			user.POST("/assign-passcodes", UserHandler.AssignPasscodes)
			user.POST("/import-users-with-emails", UserHandler.UpdateUsersImport)
		}

		{ // ==== ADMIN
			admin.Use(h.RequireAnyRole(m.RoleSystemAdminManager, m.RoleSystemSuperAdmin))
			admin.GET("", UserHandler.GetAdmins)
			admin.POST("/:id", UserHandler.MakeAdmin)
			admin.PUT("/", UserHandler.UpdateAdmin)
			admin.DELETE("/:id", UserHandler.DeleteAdmin)
			admin.POST("/add-manager/:id", h.RequireRole(m.RoleSystemSuperAdmin), UserHandler.MakeAdminManager)
		}

		{ // ==== BLOG POSTS
			blog := admin.Group("/blog")
			blog.Use(h.RequireAnyRole(m.RoleContentBlogMgr, m.RoleSystemSuperAdmin))
			blog.GET("", CmsHandler.GetAllBlogPosts)
			blog.GET("/:id", CmsHandler.GetBlogPostById)
			blog.GET("/slug/:slug", CmsHandler.GetBlogPostBySlug)
			blog.POST("", CmsHandler.CreateBlogPost)
			blog.PUT("", CmsHandler.UpdateBlogPost)
			blog.DELETE("/:id", CmsHandler.DeleteBlogPost)
		}

		{ // ==== GALLERY
			gallery := admin.Group("/gallery")
			gallery.Use(h.RequireAnyRole(m.RoleContentEditor, m.RoleSystemSuperAdmin))
			gallery.POST("", GalleryHandler.Upload)
			gallery.GET("", GalleryHandler.GetAll)
			gallery.GET("/:id", GalleryHandler.GetByID)
			gallery.DELETE("", GalleryHandler.CleanGallery)
		}

		// TODO: Add bot commands

		{ // ==== FORMS
			form := admin.Group("/form")
			form.Use(h.RequireAnyRole(m.RoleContentFormMgr, m.RoleSystemSuperAdmin))

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

			// form.GET("/user-response/:id", FormHandler.GetResponseByID)
			// form.GET("/user-responses/:id", FormHandler.GetUserResponsesForForm)

			// form.GET("/responses/:id", FormHandler.GetResponsesByFormID)
			// form.PUT("/response-status", FormHandler.UpdateResponseStatus)
			// form.DELETE("/response/:id", FormHandler.DeleteResponse)

		}

		{ // ==== TEAM MEMBERS
			team := admin.Group("/team")
			team.Use(h.RequireAnyRole(m.RoleContentEditor, m.RoleSystemSuperAdmin))
			team.POST("", CmsHandler.CreateTeamMember)
			team.GET("", CmsHandler.GetAllTeamMembers)
			team.GET("/:id", CmsHandler.GetTeamMemberByID)
			team.PUT("", CmsHandler.UpdateTeamMember)
			team.DELETE("/:id", CmsHandler.DeleteTeamMember)
		}

		{ // ==== EVENTS
			event := admin.Group("/event")
			event.Use(h.RequireAnyRole(m.RoleContentEventMgr, m.RoleSystemSuperAdmin))
			event.GET("", EventHandler.GetAllEvents)
			event.GET("/:id", EventHandler.GetEventByID)
			event.POST("", EventHandler.CreateEvent)
			event.PUT("", EventHandler.UpdateEvent)
			event.DELETE("/:id", EventHandler.DeleteEvent)
			event.GET("/send-all-emails", CertificateHandler.SendCertificatesEmailsForEvent)
			event.POST("/import-users/:id", EventHandler.ImportUsers)
		}

		{ // ==== CERTIFICATES
			certificate := admin.Group("/certificate")
			certificate.Use(h.RequireAnyRole(m.RoleCertifier, m.RoleSystemSuperAdmin))
			certificate.GET("/generate-all-for-event", CertificateHandler.MakeCertificatesForEvent)
		}

		{ // ==== MAIL
			mail := apiV1.Group("/mail")
			mail.POST("", MailHandler.SendMail)
		}
	}
	return r
}
