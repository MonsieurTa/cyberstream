package main

import (
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&entity.Credential{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&entity.PublicInfo{}); err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := gorm.Open(postgres.Open("host=psql user=docker password=secret dbname=docker port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	err = Migrate(db)
	if err != nil {
		panic("failed to migrate models")
	}

	router := gin.Default()
	validator.RegisterCustomValidations(router)

	app := NewApp(db, router)

	app.MakeHandlers()

	app.Run()
}
