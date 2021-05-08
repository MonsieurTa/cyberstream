package main

import (
	"fmt"
	"os"

	"github.com/MonsieurTa/hypertube/common/infrastructure/database"
	"github.com/MonsieurTa/hypertube/common/validator"
	a "github.com/MonsieurTa/hypertube/pkg/api/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

func initDB() database.Database {
	format := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"

	db_uri := fmt.Sprintf(format,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"))

	db := database.NewDBGORM(db_uri, &gorm.Config{})

	err := db.Migrate()
	if err != nil {
		panic("failed to migrate")
	}
	return db
}

func initEnv() {
	env := os.Getenv("HYPERTUBE_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
}

func main() {
	initEnv()

	db := initDB()

	router := gin.Default()
	router.Use(cors.Default())

	if ok := validator.RegisterCustomValidations(router); !ok {
		panic("could not register custom validations")
	}

	app, err := a.NewApp(db.DB().(*gorm.DB), router)
	if err != nil {
		panic(err)
	}

	app.MakeHandlers()
	app.Run()
}
