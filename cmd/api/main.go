package main

import (
	"os"

	"github.com/MonsieurTa/hypertube/common/db"
	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

func initEnv() {
	env := os.Getenv("HYPERTUBE_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
}

func main() {
	initEnv()

	db := db.InitDB(&db.PSQLConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Db:       os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}, &gorm.Config{})

	router := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AddExposeHeaders("Authorization")
	cfg.AddAllowHeaders("Authorization")
	cfg.AddAllowMethods("OPTIONS")
	cfg.AllowOrigins = []string{"http://localhost:8081", "https://cyberstream.digital"}
	cfg.AllowCredentials = true

	router.Use(cors.New(cfg))

	if ok := validator.RegisterCustomValidations(router); !ok {
		panic("could not register custom validations")
	}

	app, err := server.NewServer(db.DB().(*gorm.DB), router)
	if err != nil {
		panic(err)
	}

	app.MakeHandlers()
	app.Run()
}
