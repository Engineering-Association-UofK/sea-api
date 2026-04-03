package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
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
	var email models.UsersEmails
	if err := ctx.ShouldBindJSON(&email); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	if err := h.MailService.SendUsersEmails(email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, nil)
}
