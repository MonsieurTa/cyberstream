package main

import (
	"github.com/MonsieurTa/hypertube/api/handler"
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/infrastructure/repository"
	auth "github.com/MonsieurTa/hypertube/usecase/authentication"
	"github.com/MonsieurTa/hypertube/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/usecase/state"
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

	stateInMem := repository.NewStateInMem()
	stateService := state.NewService(stateInMem)

	credentialRepo := repository.NewCredentialGORM(a.db)
	authService := auth.NewService(credentialRepo)

	handler.MakeUsersHandlers(v1.Group("/users"), userService)
	handler.MakeOAuth2Handlers(v1.Group("/oauth"), ftService, stateService, authService)
}

func (a *App) Run() error {
	return a.router.Run(config.PORT)
}
