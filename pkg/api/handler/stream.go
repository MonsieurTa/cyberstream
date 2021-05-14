package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/movie"
	"github.com/gin-gonic/gin"
)

const (
	ERR_HTTP_POST            = "error_http_post"
	ERR_RESP_BODY            = "error_resp_body"
	ERR_UNMARSHALL           = "error_unmarshall"
	ERR_MARSHALL             = "error_marshall"
	ERR_INVALID_PARAMS       = "error_invalid_params"
	ERR_REPOSITORY           = "error_repository"
	ERR_SERVICE_REGISTRATION = "error_service_registration"
)

func RequestStream(movieService movie.UseCase) gin.HandlerFunc {
	MEDIA_ENDPOINT := `http://` + os.Getenv("MEDIA_HOST") + `:` + os.Getenv("MEDIA_PORT") + `/stream`
	return func(c *gin.Context) {
		validator := validator.NewUnstoredMovieValidator()
		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_INVALID_PARAMS: err.Error()})
			return
		}
		input := validator.Value()

		res, err := movieService.FindByID(input.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{ERR_REPOSITORY: err.Error()})
			return
		}

		if res != nil {
			streamResponse := entity.NewStreamResponse(res.Path)
			c.JSON(http.StatusOK, streamResponse)
			return
		}

		streamReq := entity.NewStreamRequest(input.Name, input.Magnet)
		data, err := json.Marshal(gin.H{"stream_request": streamReq})
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{ERR_MARSHALL: err.Error()})
			return
		}

		resp, err := http.Post(MEDIA_ENDPOINT, "application/json", bytes.NewReader(data))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{ERR_HTTP_POST: err.Error()})
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{ERR_RESP_BODY: err.Error()})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, string(b))
			return
		}

		var streamResponse entity.StreamResponse

		err = json.Unmarshal(b, &streamResponse)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ERR_UNMARSHALL: err.Error()})
			return
		}
		input.Path = streamResponse.Url

		_, err = movieService.Register(&input)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{ERR_SERVICE_REGISTRATION: err.Error()})
			return
		}

		c.JSON(http.StatusOK, streamResponse)
	}
}
