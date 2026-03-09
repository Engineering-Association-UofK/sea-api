package main

import (
	"sea-api/cmd/routes"
	"sea-api/internal/config"
	"sea-api/internal/handlers"
	"sea-api/internal/storage"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

type User struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func main() {
	Init()

	db := storage.NewMySQLConnection()

	routes.UserHandler = handlers.NewUserHandler(db)
	routes.EventHandler = handlers.NewEventHandler(db)

	r := routes.SetupRouter()
	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func Init() {
	gin.SetMode(gin.ReleaseMode)
	err := config.Load()
	if err != nil {
		panic(err)
	}
}
