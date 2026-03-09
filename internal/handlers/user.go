package handlers

import (
	"log/slog"
	"sea-api/internal/response"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{service: services.NewUserService(db)}
}

func (u *UserHandler) GetAll(c *gin.Context) {
	users, err := u.service.GetAll()
	if err != nil {
		slog.Error("Error fetching users from the database:"+err.Error(), "handler", "UserHandler", "function", "GetAll")
		response.InternalServerError(c)
		return
	}

	c.JSON(200, users)
}
