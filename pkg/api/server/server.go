package server

import (
	"os"

	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/pkg/api/handler"
	"github.com/MonsieurTa/hypertube/pkg/api/internal/inmem"
	"github.com/MonsieurTa/hypertube/pkg/api/middleware"
	auth "github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/jackett"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/provider"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/state"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/stream"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/user"
	"github.com/gin-gonic/gin"
	gojackett "github.com/webtor-io/go-jackett"
	"gorm.io/gorm"
)

type Server struct {
	db       *gorm.DB
	router   *gin.Engine
	services *Services
}

type Services struct {
	// tmp cache for token `state`
	auth     auth.UseCase
	fortytwo fortytwo.UseCase
	jackett  jackett.UseCase
	provider provider.UseCase
	state    state.UseCase
	stream   stream.UseCase
	user     user.UseCase
}

func NewServer(db *gorm.DB, router *gin.Engine) (*Server, error) {
	services, err := registerServices(db)
	if err != nil {
		return nil, err
	}
	return &Server{
		db,
		router,
		services,
	}, nil
}

func (s *Server) MakeHandlers() {
	secret := os.Getenv("JWT_SECRET")
	v1 := s.router.Group("/api")

	v1.POST("/stream", handler.RequestStream(s.services.stream))

	auth := v1.Group("/oauth")
	auth.POST("/token", handler.AccessTokenGeneration(s.services.auth))
	auth.GET("/fortytwo/callback", handler.RedirectCallback(s.services.fortytwo, s.services.state))
	auth.GET("/fortytwo/authorize_uri", handler.GetAuthorizeURI(s.services.fortytwo, s.services.state))

	users := v1.Group("/users")
	users.GET("/", middleware.Auth(secret))    // TODO
	users.GET("/:id", middleware.Auth(secret)) // TODO
	users.PATCH("/:id", middleware.Auth(secret), handler.UsersUpdate(s.services.user))
	users.POST("/", handler.UsersRegistration(s.services.user))

	videos := v1.Group("/videos") // TODO
	videos.GET("/")               // TODO
	videos.GET("/:id")            // TODO

	anime := v1.Group("/jackett") // TODO
	anime.GET("/search", handler.JackettSearch(s.services.jackett))
	anime.GET("/categories", handler.JackettCategories(s.services.jackett))
	anime.GET("/indexers", handler.JackettIndexers(s.services.jackett))
}

func registerServices(db *gorm.DB) (*Services, error) {
	ftService, err := fortytwo.NewService()
	if err != nil {
		return nil, err
	}

	jackettRepo := gojackett.NewJackett(&gojackett.Settings{
		ApiURL: os.Getenv("JACKETT_API_URL"),
		ApiKey: os.Getenv("JACKETT_API_KEY"),
	})

	providerRepo := repository.NewProviderGORM(db)
	videoRepository := repository.NewVideoGORM(db)
	stateInMem := inmem.NewStateInMem()
	userRepo := repository.NewUserGORM(db)

	authService := auth.NewService(userRepo)
	jackettService := jackett.NewService(jackettRepo)
	providerService, err := provider.NewService(providerRepo)
	if err != nil {
		return nil, err
	}

	stateService := state.NewService(stateInMem)
	streamService := stream.NewService(videoRepository)
	userService := user.NewService(userRepo)
	return &Services{
		auth:     authService,
		fortytwo: ftService,
		jackett:  jackettService,
		provider: providerService,
		state:    stateService,
		stream:   streamService,
		user:     userService,
	}, nil
}

func (s *Server) Run() error {
	port := ":" + os.Getenv("API_PORT")
	return s.router.Run(port)
}
