package handlers

import (
	"net/http"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	AccountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		AccountService: accountService,
	}
}

func (a *AccountHandler) GetProfile(c *gin.Context) {
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	profile, err := a.AccountService.GetProfile(c.Request.Context(), claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, profile)
}

func (a *AccountHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	err := a.AccountService.UpdateProfile(claims, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Profile updated successfully", req.ID, c)
}

func (a *AccountHandler) UpdatePicture(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024)

	file, _, err := c.Request.FormFile("picture")
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	err = a.AccountService.UpdateProfilePicture(c.Request.Context(), claims, file)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Profile picture updated successfully", claims.UserID, c)
}

func (a *AccountHandler) UpdatePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}
	err := a.AccountService.UpdatePassword(claims, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.NewTransactionResponse(200, "Password updated successfully", claims.UserID, c)
}

func (a *AccountHandler) UpdateEmail(c *gin.Context) {
	var req models.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	err := a.AccountService.UpdateEmail(claims, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.NewTransactionResponse(200, "Email updated successfully", claims.UserID, c)
}

func (a *AccountHandler) UpdateUsername(c *gin.Context) {
	var req models.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}
	err := a.AccountService.UpdateUsername(claims, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.NewTransactionResponse(200, "Username updated successfully", claims.UserID, c)
}

func (a *AccountHandler) CheckUsernameAvailability(c *gin.Context) {
	var req models.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	available := a.AccountService.IsUsernameAvailable(req)
	c.JSON(http.StatusOK, models.CheckUsername{Available: available})
}
