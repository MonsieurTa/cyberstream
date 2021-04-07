package main

import (
	"github.com/MonsieurTa/hypertube/api/handler"
	"github.com/MonsieurTa/hypertube/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/usecase/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewApp(db *gorm.DB, router *gin.Engine) *App {
	return &App{
		db,
		router,
	}
}

func (a *App) MakeHandlers() {
	v1 := a.router.Group("/api")

	userRepo := repository.NewUserGORM(a.db)
	userService := user.NewService(userRepo)

	ftService, _ := fortytwo.NewService()

	handler.MakeUsersHandlers(v1.Group("/users"), userService)
	handler.MakeFortyTwoAuthHandlers(v1.Group("/auth"), ftService)
}

func (a *App) Run() error {
	return a.router.Run()
}
