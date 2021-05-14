package app

import (
	"os"

	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/pkg/api/handler"
	"github.com/MonsieurTa/hypertube/pkg/api/internal/inmem"
	s "github.com/MonsieurTa/hypertube/pkg/api/internal/subsplease"
	"github.com/MonsieurTa/hypertube/pkg/api/middleware"
	auth "github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/movie"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/provider"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/state"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/subsplease"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	db       *gorm.DB
	router   *gin.Engine
	services *Services
}

type Services struct {
	// tmp cache for token `state`
	state state.UseCase

	auth     auth.UseCase
	fortytwo fortytwo.UseCase

	user     user.UseCase
	provider provider.UseCase
	movie    movie.UseCase

	subsplease subsplease.UseCase
}

func NewApp(db *gorm.DB, router *gin.Engine) (*App, error) {
	services, err := registerServices(db)
	if err != nil {
		return nil, err
	}
	return &App{
		db,
		router,
		services,
	}, nil
}

func (a *App) MakeHandlers() {
	secret := os.Getenv("JWT_SECRET")
	v1 := a.router.Group("/api")

	v1.POST("/stream", handler.RequestStream(a.services.movie))

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

	subsplease := v1.Group("/subsplease")
	subsplease.GET("/latest", handler.SubsPleaseLatestEpisodes(a.services.subsplease))
}

func registerServices(db *gorm.DB) (*Services, error) {
	userRepo := repository.NewUserGORM(db)
	providerRepo := repository.NewProviderGORM(db)
	stateInMem := inmem.NewStateInMem()
	subspleaseRepo := s.NewSubsPlease()
	movieRepository := repository.NewMovieGORM(db)

	ftService, err := fortytwo.NewService()
	if err != nil {
		return nil, err
	}

	stateService := state.NewService(stateInMem)
	authService := auth.NewService(userRepo)
	userService := user.NewService(userRepo)
	providerService, err := provider.NewService(providerRepo)
	movieService := movie.NewService(movieRepository)

	subspleaseService := subsplease.NewService(subspleaseRepo)
	if err != nil {
		return nil, err
	}
	return &Services{
		stateService,
		authService,
		ftService,
		userService,
		providerService,
		movieService,
		subspleaseService,
	}, nil
}

func (a *App) Run() error {
	port := ":" + os.Getenv("API_PORT")
	return a.router.Run(port)
}
