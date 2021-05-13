package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/pkg/api/usecase/subsplease"
	"github.com/gin-gonic/gin"
)

const (
	ERR_SUBSPLEASE_UNAVAILABLE = "err_subsplease_unavailable"
)

func SubsPleaseLatestEpisodes(service subsplease.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		episodes, err := service.Latest()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{ERR_SUBSPLEASE_UNAVAILABLE: err.Error()})
			return
		}
		c.JSON(http.StatusOK, episodes)
	}
}
