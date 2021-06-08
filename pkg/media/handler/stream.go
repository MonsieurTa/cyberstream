package handler

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/file"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/MonsieurTa/hypertube/pkg/media/usecase/media"
	"github.com/anacrolix/torrent"
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

		streamResp := service.StreamMagnet(streamReq.Magnet)

		c.JSON(http.StatusOK, &streamResp)
	}
}

type ContentRequest struct {
	Name string `form:"name" json:"name"`
	Hash string `form:"hash" json:"hash"`
}

// mp4
func ServeContent(service media.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cr ContentRequest

		err := c.ShouldBindQuery(&cr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_params": err.Error()})
			return
		}

		relpath, _ := url.PathUnescape(cr.Name)
		filePath := os.Getenv("STATIC_FILES_PATH") + "/" + relpath
		fmt.Println(filePath)
		fileExists := file.Exists(filePath)

		t, torrentExists := getTorrentFromHash(service, cr.Hash)

		if !fileExists && !torrentExists {
			c.String(http.StatusNotFound, "not found")
			return
		}

		var rs io.ReadSeeker

		if fileExists && !torrentExists {
			rs, _ = os.Open(filePath)
		} else {
			rs = t.NewReader()
		}

		w := c.Writer
		req := c.Request
		http.ServeContent(w, req, cr.Name, time.Time{}, rs)
	}
}

func getTorrentFromHash(service media.UseCase, hash string) (*torrent.Torrent, bool) {
	var hexHash [20]byte

	hex, err := hex.DecodeString(hash)
	if err != nil {
		return nil, false
	}
	copy(hexHash[:], hex[:20])
	return service.Torrent(hexHash)
}
