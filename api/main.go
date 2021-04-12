package main

import (
	"fmt"

	a "github.com/MonsieurTa/hypertube/api/app"
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/infrastructure/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(model.UserModelGORM()); err != nil {
		return err
	}
	if err := db.AutoMigrate(model.CredentialsModelGORM()); err != nil {
		return err
	}
	if err := db.AutoMigrate(model.PublicInfoModelGORM()); err != nil {
		return err
	}
	return nil
}

func main() {
	format := "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable"
	db_uri := fmt.Sprintf(format, config.POSTGRES_HOST, config.POSTGRES_USER, config.POSTGRES_PASSWORD, config.POSTGRES_DB)

	db, err := gorm.Open(postgres.Open(db_uri), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	err = Migrate(db)
	if err != nil {
		panic("failed to migrate models")
	}

	router := gin.Default()
	if ok := validator.RegisterCustomValidations(router); !ok {
		panic("could not register custom validations")
	}

	app, err := a.NewApp(db, router)
	if err != nil {
		panic(err)
	}

	app.MakeHandlers()
	app.Run()
}
