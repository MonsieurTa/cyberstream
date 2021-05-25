package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/MonsieurTa/hypertube/pkg/media/usecase/media"
	"github.com/gin-gonic/gin"
)

const (
	ERR_MAGNET         = "error_magnet"
	ERR_TORRENT_CLIENT = "error_torrent_client"
	ERR_VALIDATION     = "error_validation"
)

func Stream(service media.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var streamReq entity.StreamRequest

		err := common.Bind(c, &streamReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_params": err.Error()})
			return
		}

		streamResp, err := service.StreamMagnet(streamReq.Magnet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_MAGNET: err.Error()})
			return
		}

		c.JSON(http.StatusOK, &streamResp)
	}
}
