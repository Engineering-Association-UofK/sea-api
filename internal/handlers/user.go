package handlers

import (
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) GetAll(c *gin.Context) {
	users, err := u.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, users)
}
