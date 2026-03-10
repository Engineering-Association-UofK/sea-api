package handlers

import (
	"log/slog"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type MailHandler struct {
	MailService *services.MailService
}

func NewMailHandler(db *sqlx.DB) *MailHandler {
	return &MailHandler{
		MailService: services.NewMailService(
			services.NewUserService(db),
		),
	}
}

func (h *MailHandler) SendMail(ctx *gin.Context) {
	var email models.UserEmails
	if err := ctx.ShouldBindJSON(&email); err != nil {
		response.BadRequest(ctx)
		return
	}

	if err := h.MailService.SendUserEmails(email); err != nil {
		slog.Error("error sending email", "error", err)
		response.InternalServerError(ctx)
		return
	}

	ctx.JSON(200, gin.H{"message": "Email sent successfully"})
}
