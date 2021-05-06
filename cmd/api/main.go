package main

import (
	"fmt"

	"github.com/MonsieurTa/hypertube/common/infrastructure/database"
	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/config"
	a "github.com/MonsieurTa/hypertube/server-api/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initDB() database.Database {
	format := "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable"
	db_uri := fmt.Sprintf(format, config.POSTGRES_HOST, config.POSTGRES_USER, config.POSTGRES_PASSWORD, config.POSTGRES_DB)
	db := database.NewDBGORM(db_uri, &gorm.Config{})

	err := db.Migrate()
	if err != nil {
		panic("failed to migrate")
	}
	return db
}

func main() {
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
