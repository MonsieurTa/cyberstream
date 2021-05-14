package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/stream"
	"github.com/gin-gonic/gin"
)

const (
	ERR_INVALID_PARAMS       = "error_invalid_params"
	ERR_REPOSITORY           = "error_repository"
	ERR_SERVICE_REGISTRATION = "error_service_registration"
	ERR_STREAM_SERVICE       = "error_stream_service"
)

func RequestStream(streamService stream.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUnstoredMovieValidator()
		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_INVALID_PARAMS: err.Error()})
			return
		}
		input := validator.Value()

		url, err := streamService.Stream(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_STREAM_SERVICE: err.Error()})
			return
		}
		c.JSON(http.StatusOK, url)
	}
}
