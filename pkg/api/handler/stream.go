package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/gin-gonic/gin"
)

const (
	ERR_HTTP_POST  = "error_http_post"
	ERR_RESP_BODY  = "error_resp_body"
	ERR_UNMARSHALL = "error_unmarshall"
)

var MEDIA_ENDPOINT = `http://` + os.Getenv("MEDIA_HOST") + `:` + os.Getenv("MEDIA_PORT") + `/stream`

func RequestStream(c *gin.Context) {
	resp, err := http.Post(MEDIA_ENDPOINT, "application/json", c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_HTTP_POST: err.Error(),
		})
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_RESP_BODY: err.Error(),
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, string(b))
		return
	}

	var streamResponse entity.StreamResponse

	err = json.Unmarshal(b, &streamResponse)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_UNMARSHALL: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, streamResponse)
}
