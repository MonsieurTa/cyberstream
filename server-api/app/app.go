package app

import (
	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/server-api/handler"
	"github.com/MonsieurTa/hypertube/server-api/internal/inmem"
	"github.com/MonsieurTa/hypertube/server-api/middleware"
	auth "github.com/MonsieurTa/hypertube/server-api/usecase/authentication"
	"github.com/MonsieurTa/hypertube/server-api/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/server-api/usecase/provider"
	"github.com/MonsieurTa/hypertube/server-api/usecase/state"
	"github.com/MonsieurTa/hypertube/server-api/usecase/user"
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

	user     user.UseCase
	provider provider.UseCase
}

func NewApp(db *gorm.DB, router *gin.Engine) (*App, error) {
	userRepo := repository.NewUserGORM(db)
	providerRepo := repository.NewProviderGORM(db)
	stateInMem := inmem.NewStateInMem()

	ftService, err := fortytwo.NewService()
	if err != nil {
		return nil, err
	}

	stateService := state.NewService(stateInMem)
	authService := auth.NewService(userRepo)
	userService := user.NewService(userRepo)
	providerService, err := provider.NewService(providerRepo)
	if err != nil {
		return nil, err
	}
	return &App{
		db,
		router,
		Services{
			stateService,
			authService,
			ftService,
			userService,
			providerService,
		},
	}, nil
}

func (a *App) MakeHandlers() {
	v1 := a.router.Group("/api")

	v1.POST("/stream", handler.RequestStream)

	auth := v1.Group("/oauth")
	auth.POST("/token", handler.AccessTokenGeneration(a.services.auth))
	auth.GET("/fortytwo/callback", handler.RedirectCallback(a.services.fortytwo, a.services.state))
	auth.GET("/fortytwo/authorize_uri", handler.GetAuthorizeURI(a.services.fortytwo, a.services.state))

	users := v1.Group("/users")
	users.GET("/", middleware.Auth(config.JWT_SECRET))    // TODO
	users.GET("/:id", middleware.Auth(config.JWT_SECRET)) // TODO
	users.PATCH("/:id", middleware.Auth(config.JWT_SECRET), handler.UsersUpdate(a.services.user))
	users.POST("/", handler.UsersRegistration(a.services.user))

	movies := v1.Group("/movies") // TODO
	movies.GET("/")               // TODO
	movies.GET("/:id")            // TODO
}

func (a *App) Run() error {
	return a.router.Run(config.PORT)
}
