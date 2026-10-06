package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/authmodels"
	"sea-api/internal/models/usermodels"
	"sea-api/internal/response"
	"sea-api/internal/services/userservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *userservice.UserService
}

func NewUserHandler(service *userservice.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ======== GET ALL ========

// GetAll godocs
//
//	@Summary		Get all users
//	@Description	Get a list of all users with pagination
//	@Tags			User::Management
//	@Produce		json
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Success		200		{object}	usermodels.UserListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/user/all [get]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) GetAll(c *gin.Context) {
	var req models.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))

		return
	}
	resp, err := u.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}

// GetAdmins godocs
//
//	@Summary		Get all admins
//	@Description	Get a list of all administrative users
//	@Tags			User::Admin
//	@Produce		json
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Success		200	{object}		usermodels.AdminResponseList
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin [get]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) GetAdmins(c *gin.Context) {
	req := &models.ListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	ctx := c.Request.Context()
	admins, err := u.service.GetAdmins(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, admins)
}

// ======== GET ========

// GetByID godocs
//
//	@Summary		Get user by ID
//	@Description	Get user details by their ID
//	@Tags			User::Management
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	usermodels.UserResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/user/{id} [get]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	user, err := u.service.GetByUserID(c.Request.Context(), intId)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(200, user)
}

// GetByUsername godocs
//
//	@Summary		Get user by username
//	@Description	Get user details by their username
//	@Tags			User::Management
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200			{object}	usermodels.UserResponse
//	@Failure		400			{object}	response.BaseError
//	@Failure		401			{object}	response.BaseError
//	@Failure		404			{object}	response.BaseError
//	@Failure		500			{object}	response.BaseError
//	@Router			/admin/user/username/{username} [get]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := u.service.GetByUsername(c.Request.Context(), username)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(200, user)
}

// GetPasscode godocs
//
//	@Summary		Get Passcode
//	@Description	Get the registration passcode for user by their ID
//	@Tags			User::Management
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	usermodels.GetPasscodeResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/user/passcode/{id} [get]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) GetPasscode(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	passcode, err := u.service.GetPasscode(intId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, passcode)
}

// ======== MAKE CHANGES ========

// Update godocs
//
//	@Summary		Update user
//	@Description	Update user profile details by administration
//	@Tags			User::Management
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usermodels.UpdateProfileRequest	true	"Update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/user [put]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) Update(c *gin.Context) {
	var req usermodels.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err := u.service.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Profile updated successfully", req.ID, c)
}

// MakeAdmin godocs
//
//	@Summary		Make user admin
//	@Description	Assign administrative roles to a user
//	@Tags			User::Admin
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/{id} [post]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) MakeAdmin(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err = u.service.AddAdmin(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User added as admin successfully", intId, c)
}

// MakeAdminManager godocs
//
//	@Summary		Make user admin manager
//	@Description	Assign admin manager role to a user
//	@Tags			User::Admin
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/add-manager/{id} [post]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) MakeAdminManager(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err = u.service.MakeAdminManager(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User added as admin manager successfully", intId, c)
}

// RemoveAdminManager godocs
//
//	@Summary		Remove admin manager
//	@Description	Remove admin manager role from a user
//	@Tags			User::Admin
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/remove-manager/{id} [delete]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) RemoveAdminManager(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err = u.service.RemoveAdminManager(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User removed as admin manager successfully", intId, c)
}

// UpdateAdmin godocs
//
//	@Summary		Update admin roles
//	@Description	Update administrative roles for an existing admin
//	@Tags			User::Admin
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usermodels.AdminRequest	true	"Admin update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin [put]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) UpdateAdmin(c *gin.Context) {
	var req usermodels.AdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err := u.service.UpdateAdminRoles(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Admin updated successfully", req.ID, c)
}

// DeleteAdmin godocs
//
//	@Summary		Delete admin
//	@Description	Remove administrative roles from a user
//	@Tags			User::Admin
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/{id} [delete]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	err = u.service.RemoveAdmin(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Admin removed successfully", intId, c)
}

// ======== SPECIAL ========

// Suspend godocs
//
//	@Summary		Suspend user
//	@Description	Suspend a user account for a specified duration
//	@Tags			User::Management
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.SuspensionRequest	true	"Suspension data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/user/suspend [post]
//
//	@Security		ApiKeyAuth
func (u *UserHandler) Suspend(c *gin.Context) {
	var req models.SuspensionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*authmodels.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	err := u.service.Suspend(&req, claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User suspended successfully", req.UserID, c)
}
