package handlers

import (
	"encoding/json"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
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

// CheckState godocs
//
//	@Summary		Check registration step
//	@Description	Check which registration step the user is in right now.
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.CheckRegistrationRequest	true	"Request body"
//
//	@Success		201	{object}	response.CheckRegistrationResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/register/check [post]
func (h *AuthHandler) CheckState(c *gin.Context) {
	var req models.CheckRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	res, _, err := h.AuthService.CheckRegistration(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}

// DoRegistrationStep godocs
//
//	@Summary		Register Step
//	@Description	Do registration step with it's data and step number
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.CheckRegistrationRequest	true	"Request body"
//
//	@Success		201	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/register/step [post]
func (h *AuthHandler) DoRegistrationStep(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	slog.Debug("Request started", "Step", req.Step)

	var err error = nil

	switch req.Step {
	case 0:
		var data models.InitialRegistrationRequest
		if err = json.Unmarshal(req.Data, &data); err == nil {
			err = h.AuthService.InitialRegistration(&data)
		}
	case 1:
		var data models.PasswordRegistrationRequest
		if err = json.Unmarshal(req.Data, &data); err == nil {
			err = h.AuthService.CredentialsRegistration(&data)
		}
	case 2:
		var data models.DetailsRegistrationRequest
		if err = json.Unmarshal(req.Data, &data); err == nil {
			err = h.AuthService.DetailsRegistration(&data)
		}
	case 3:
		var data models.UsernameRegistrationRequest
		if err = json.Unmarshal(req.Data, &data); err == nil {
			err = h.AuthService.UsernameRegistration(&data)
		}
	case 5:
		var data models.PasswordRegistrationRequest
		if err = json.Unmarshal(req.Data, &data); err == nil {
			err = h.AuthService.CredentialsRegistration(&data)
		}
	default:
		c.Error(errs.New(errs.BadRequest, "Provided step number is out of range", nil))
		return
	}

	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Step registered successfully", 0, c)
}

// ForgotPassword godocs
//
//	@Summary		Forgot Password
//	@Description	Forgot password endpoint to send a password reset email
//	@Tags			Auth
//	@Produce		json
//	@Param			body	body	models.ForgotPasswordRequest	true	"Request body"
//
//	@Success		201	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.AuthService.ForgotPassword(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Password reset email sent", 0, c)
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
