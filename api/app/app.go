package app

import (
	"github.com/MonsieurTa/hypertube/api/handler"
	"github.com/MonsieurTa/hypertube/api/middleware"
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
	db       *gorm.DB
	router   *gin.Engine
	services Services
}

type Services struct {
	// tmp cache for token `state`
	state state.UseCase

	auth     auth.UseCase
	fortytwo fortytwo.UseCase

	user user.UseCase
}

func NewApp(db *gorm.DB, router *gin.Engine) (*App, error) {
	userRepo := repository.NewUserGORM(db)
	stateInMem := repository.NewStateInMem()

	ftService, err := fortytwo.NewService()
	if err != nil {
		return nil, err
	}

	stateService := state.NewService(stateInMem)
	authService := auth.NewService(userRepo)
	userService := user.NewService(userRepo)
	return &App{
		db,
		router,
		Services{
			stateService,
			authService,
			ftService,
			userService,
		},
	}, nil
}

func (a *App) MakeHandlers() {
	v1 := a.router.Group("/api")

	a.makeAuthHandlers(v1.Group("/oauth"))
	a.makeUsersHandlers(v1.Group("/users"))
}

func (a *App) makeAuthHandlers(g *gin.RouterGroup) {
	g.POST("/token", handler.AccessTokenGeneration(a.services.auth))

	g.GET("/fortytwo/callback", handler.RedirectCallback(a.services.fortytwo, a.services.state))
	g.GET("/fortytwo/authorize_uri", handler.GetAuthorizeURI(a.services.fortytwo, a.services.state))
}

func (a *App) makeUsersHandlers(g *gin.RouterGroup) {
	g.POST("/", handler.UsersRegistration(a.services.user))
	g.PATCH("/", middleware.Auth(config.JWT_SECRET), handler.UsersUpdate(a.services.user))
}

func (a *App) Run() error {
	return a.router.Run(config.PORT)
}
