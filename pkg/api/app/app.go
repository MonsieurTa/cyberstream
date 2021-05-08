package app

import (
	"os"

	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/pkg/api/handler"
	"github.com/MonsieurTa/hypertube/pkg/api/internal/inmem"
	"github.com/MonsieurTa/hypertube/pkg/api/middleware"
	auth "github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/provider"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/state"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/user"
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
	secret := os.Getenv("JWT_SECRET")
	v1 := a.router.Group("/api")

	v1.POST("/stream", handler.RequestStream)

	auth := v1.Group("/oauth")
	auth.POST("/token", handler.AccessTokenGeneration(a.services.auth))
	auth.GET("/fortytwo/callback", handler.RedirectCallback(a.services.fortytwo, a.services.state))
	auth.GET("/fortytwo/authorize_uri", handler.GetAuthorizeURI(a.services.fortytwo, a.services.state))

	users := v1.Group("/users")
	users.GET("/", middleware.Auth(secret))    // TODO
	users.GET("/:id", middleware.Auth(secret)) // TODO
	users.PATCH("/:id", middleware.Auth(secret), handler.UsersUpdate(a.services.user))
	users.POST("/", handler.UsersRegistration(a.services.user))

	movies := v1.Group("/movies") // TODO
	movies.GET("/")               // TODO
	movies.GET("/:id")            // TODO
}

func (a *App) Run() error {
	port := ":" + os.Getenv("API_PORT")
	return a.router.Run(port)
}
