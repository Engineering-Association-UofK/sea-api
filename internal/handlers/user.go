package handlers

import (
	"io"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ======== GET ALL ========

func (u *UserHandler) GetAll(c *gin.Context) {
	var req models.UserListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	resp, err := u.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}
func (u *UserHandler) GetAllTempUsers(c *gin.Context) {
	var req models.UserListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	resp, err := u.service.GetAllTempUsers(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}

func (u *UserHandler) GetAdmins(c *gin.Context) {
	admins, err := u.service.GetAdmins()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, admins)
}

// ======== GET ========

func (u *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	user, err := u.service.GetByUserID(c.Request.Context(), intId)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(200, user)
}

func (u *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := u.service.GetByUsername(c.Request.Context(), username)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(200, user)
}

func (u *UserHandler) GetTempUserPasscode(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	passcode, err := u.service.GetTempUserPasscode(intId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, passcode)
}

// ======== MAKE CHANGES ========

func (u *UserHandler) Update(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	err := u.service.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Profile updated successfully", req.ID, c)
}

func (u *UserHandler) MakeAdmin(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}
	err = u.service.AddAdmin(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User added as admin successfully", intId, c)
}

func (u *UserHandler) MakeAdminManager(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}
	err = u.service.MakeAdminManager(intId)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User added as admin manager successfully", intId, c)
}

func (u *UserHandler) UpdateAdmin(c *gin.Context) {
	var req models.AdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	err := u.service.UpdateAdminRoles(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Admin updated successfully", req.ID, c)
}

func (u *UserHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(c)
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

func (u *UserHandler) Suspend(c *gin.Context) {
	var req models.SuspensionRequest
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

	err := u.service.Suspend(&req, claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User suspended successfully", req.UserID, c)
}

func (u *UserHandler) AssignPasscodes(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	progressChan := make(chan string)

	go u.service.AssignPasscodes(progressChan)

	c.Stream(func(w io.Writer) bool {
		msg, ok := <-progressChan
		if !ok {
			return false
		}

		c.SSEvent("message", msg)
		return true
	})

	c.JSON(200, gin.H{"message": "Passcodes assigned successfully"})
}

func (u *UserHandler) UpdateUsersImport(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = u.service.UpdateUsersImport(file)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "User updated successfully", 0, c)
}
