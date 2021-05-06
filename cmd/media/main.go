package main

import (
	"os"

	"github.com/MonsieurTa/hypertube/server-media/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	router := gin.Default()
	router.Static("/static", os.Getenv("STATIC_FILES_PATH"))

	router.POST("/stream", handler.Stream)

	router.Run(":" + os.Getenv("STATIC_FILES_PORT"))
}
