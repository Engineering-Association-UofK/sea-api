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

//	@Param			id	path	int	true	"Notification ID"

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		AccountService: accountService,
	}
}

// GetProfileSummary godocs
//
//	@Summary		Get Profile Summary
//	@Description	Get profile summary of requesting user
//	@Tags			Account:profile
//	@Produce		json
//	@Success		200	{object}	models.UserProfileSummaryResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/summary [get]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) GetProfileSummary(c *gin.Context) {
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	profile, err := a.AccountService.GetProfileSummary(c.Request.Context(), claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, profile)
}

// GetProfile godocs
//
//	@Summary		Get Profile
//	@Description	Get all profile details of requesting user
//	@Tags			Account:profile
//	@Produce		json
//	@Success		200	{object}	models.UserProfileResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account [get]
//
//	@Security		ApiKeyAuth
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

// GetCertificates godocs
//
//	@Summary		Get user certificates
//	@Description	Get all certificates associated with the requesting user
//	@Tags			Account:profile
//	@Produce		json
//	@Param			limit	query	int	true	"Content count limit"
//	@Param			page	query	int	true	"Page number"
//	@Success		200	{array}	models.CertificateListResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/certificates [get]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) GetCertificates(c *gin.Context) {
	var req models.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request, need limit number", nil))

		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	certs, err := a.AccountService.GetCertificates(claims, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, certs)
}

// UpdateProfile godocs
//
//	@Summary		Update profile
//	@Description	Update profile text details
//	@Tags			Account:profile
//	@Produce		json
//	@Param			body	body	models.UpdateProfileRequest 	true	"Request body"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account [put]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
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

// UpdatePicture godocs
//
//	@Summary		Update profile picture
//	@Description	Update user profile picture
//	@Tags			Account:profile
//	@Produce		json
//	@Param			picture	formData	file 	true	"Upload File"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/picture [put]
//
//	@Security		ApiKeyAuth
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

// UpdatePassword godocs
//
//	@Summary		Update password
//	@Description	Update user password
//	@Tags			Account:profile
//	@Produce		json
//	@Param			body	body	models.UpdatePasswordRequest 	true	"Request body"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/password [put]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) UpdatePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
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

// UpdateEmail godocs
//
//	@Summary		Update Email
//	@Description	Update user email address
//	@Tags			Account:profile
//	@Produce		json
//	@Param			body	body	models.UpdateEmailRequest 	true	"Request body"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/email [put]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) UpdateEmail(c *gin.Context) {
	var req models.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
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

// UpdateUsername godocs
//
//	@Summary		Update Username
//	@Description	Update username
//	@Tags			Account:profile
//	@Produce		json
//	@Param			body	body	models.UpdateUsernameRequest 	true	"Request body"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/username [put]
//
//	@Security		ApiKeyAuth
func (a *AccountHandler) UpdateUsername(c *gin.Context) {
	var req models.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
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

// CheckUsernameAvailability godocs
//
//	@Summary		Update Username
//	@Description	Update username
//	@Tags			Account:profile
//	@Produce		json
//	@Param			body	body	models.UpdateUsernameRequest 	true	"Request body"
//	@Success		200	{object}	models.CheckUsername
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/auth/check-username [post]
func (a *AccountHandler) CheckUsernameAvailability(c *gin.Context) {
	var req models.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	available, err := a.AccountService.IsUsernameAvailable(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.CheckUsername{Available: available})
}
