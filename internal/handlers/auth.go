package handlers

import (
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	resp, err := h.AuthService.Login(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	err := h.AuthService.Register(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "User registered successfully", req.UserID, c)
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req models.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	err := h.AuthService.Verify(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User verified successfully", req.UserID, c)
}

func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	err := h.AuthService.SendVerificationCode(req.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Verification code sent successfully", req.UserID, c)
}
