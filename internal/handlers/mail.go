package handlers

import (
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	MailService *services.MailService
}

func NewMailHandler(mailService *services.MailService) *MailHandler {
	return &MailHandler{MailService: mailService}
}

func (h *MailHandler) SendMail(ctx *gin.Context) {
	var email models.UserEmails
	if err := ctx.ShouldBindJSON(&email); err != nil {
		response.BadRequest(ctx)
		return
	}

	if err := h.MailService.SendUserEmails(email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Email sent successfully"})
}
