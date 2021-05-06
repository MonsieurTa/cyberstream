package main

import (
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/server-media/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/static", config.STATIC_FILES_PATH)

	router.POST("/stream", handler.Stream)

	router.Run()
}
