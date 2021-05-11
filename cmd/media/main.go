package main

import (
	"os"

	"github.com/MonsieurTa/hypertube/pkg/media/handler"
	"github.com/MonsieurTa/hypertube/pkg/media/usecase/media"
	"github.com/anacrolix/torrent"
	"github.com/gin-contrib/cors"
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

func newTorrentClient() *torrent.Client {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = os.Getenv("DOWNLOAD_FILES_PATH")
	tc, err := torrent.NewClient(cfg)
	if err != nil {
		panic("could not create torrent client")
	}
	return tc
}

func main() {
	initEnv()

	tc := newTorrentClient()
	defer tc.Close()

	router := gin.Default()
	router.Use(cors.Default())

	router.Static("/", os.Getenv("STATIC_FILES_PATH"))

	mediaService := media.NewService(tc)

	router.POST("/stream", handler.Stream(mediaService))

	router.Run(":" + os.Getenv("MEDIA_PORT"))
}
