package handler

import (
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/validator"
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
		v := validator.NewStreamRequestValidator()

		err := v.Validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				ERR_VALIDATION: err.Error(),
			})
			return
		}

		magnet := v.Value().Magnet

		filepath, err := service.StreamMagnet(magnet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_MAGNET: err.Error()})
			return
		}

		resp := entity.StreamResponse{
			Url: "http://localhost" + ":" + os.Getenv("MEDIA_PORT") + "/" + filepath,
		}
		c.JSON(http.StatusOK, &resp)
	}
}
