package handlers

import (
	"sea-api/internal/errs"
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

// SendMail godocs
//
//	@Summary		Send Mail
//	@Description	Send email to multiple users
//	@Tags			Mail
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.UsersEmails	true	"Email data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/mail [post]
//
//	@Security		ApiKeyAuth
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

	response.NewTransactionResponse(200, "Email/s sent successfully", 0, ctx)
}
