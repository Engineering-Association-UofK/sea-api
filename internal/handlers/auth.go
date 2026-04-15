package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Login godocs
//
//	@Summary		Login
//	@Description	Login user
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.LoginRequest	true	"Request body"
//
//	@Success		200	{object}	models.LoginResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	resp, err := h.AuthService.Login(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}

// Register godocs
//
//	@Summary		Register
//	@Description	Register user
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.RegisterRequest	true	"Request body"
//
//	@Success		201	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.AuthService.Register(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "User registered successfully", req.UserID, c)
}

// Verify godocs
//
//		@Summary		Verify
//		@Description	Verify user
//		@Tags			Auth
//		@Produce		json
//		@Param			body	body	models.VerifyRequest	true	"Request body"
//
//		@Success		200	{object}	response.TransactionResponse
//		@Failure		400	{object}	response.BaseError
//	 @Router 			/auth/verify [post]
func (h *AuthHandler) Verify(c *gin.Context) {
	var req models.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.AuthService.Verify(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User verified successfully", req.UserID, c)
}

// SendVerificationCode godocs
//
//	@Summary		Send Verification Code
//	@Description	Send verification code to user email
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.VerifyEmailRequest	true	"Request body"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/auth/send-verification-code [post]
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.AuthService.SendVerificationCode(req.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Verification code sent successfully", req.UserID, c)
}
