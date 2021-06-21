package handler

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/file"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/MonsieurTa/hypertube/pkg/media/internal/vtt"
	"github.com/MonsieurTa/hypertube/pkg/media/usecase/iostream"
	"github.com/MonsieurTa/hypertube/pkg/media/usecase/torrenter"
	t "github.com/MonsieurTa/hypertube/pkg/media/usecase/transcoder"
	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
)

const (
	ERR_MAGNET         = "error_magnet"
	ERR_TORRENT_CLIENT = "error_torrent_client"
	ERR_VALIDATION     = "error_validation"
)

func Stream(torrenter torrenter.UseCase, transcoder t.UseCase, iostream iostream.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var streamReq entity.StreamRequest
		var streamResp entity.StreamResponse
		baseURL := "http://" + os.Getenv("MEDIA_HOST") + ":" + os.Getenv("MEDIA_PORT")

		err := common.Bind(c, &streamReq)
		if err != nil {
			streamResp.Error = err.Error()
			c.JSON(http.StatusBadRequest, &streamResp)
			return
		}

		streamInfo, err := torrenter.DownloadMagnet(&streamReq)
		if err != nil {
			streamResp.Error = err.Error()
			c.JSON(http.StatusBadRequest, &streamResp)
			return
		}

		streamResp.Name = streamInfo.StreamFile.Path()
		streamResp.InfoHash = streamInfo.InfoHash
		streamResp.Ext = streamInfo.StreamFile.Ext()
		if streamResp.Ext == ".mkv" {
			params := t.TranscoderParams{
				Reader:      streamInfo.StreamFile.NewReader(),
				FileSize:    streamInfo.StreamFile.Length(),
				PieceLength: streamInfo.StreamFile.Torrent().Info().PieceLength,
				DirName:     streamResp.InfoHash,
			}
			err := transcoder.Transcode(&params)
			if err != nil {
				streamResp.Error = err.Error()
				c.JSON(http.StatusBadRequest, &streamResp)
				return
			}

			iostream.WaitMasterPlaylist(streamResp.InfoHash, "master.m3u8")

			streamResp.MediaURL = baseURL + "/static/hls/" + streamResp.InfoHash + "/master.m3u8"
		} else {
			streamResp.MediaURL = baseURL + "/content?hash=" + streamResp.InfoHash + "&name=" + streamResp.Name
			if len(streamInfo.SubtitlesFiles) > 0 {
				converter := vtt.NewVTTConverter(os.Getenv("STATIC_FILES_PATH"), streamInfo.UnwrapSubtitlesFiles())
				filepaths := converter.Convert()
				streamResp.SubtitlesURLs = make([]string, len(streamInfo.SubtitlesFiles))
				for i, filepath := range filepaths {
					streamResp.SubtitlesURLs[i] = baseURL + "/static/" + filepath
				}
			}
		}
		c.JSON(http.StatusOK, &streamResp)
	}
}

type ContentRequest struct {
	Name string `form:"name" json:"name"`
	Hash string `form:"hash" json:"hash"`
}

// mp4
func ServeContent(service torrenter.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cr ContentRequest

		err := c.ShouldBindQuery(&cr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_params": err.Error()})
			return
		}

		relpath, _ := url.PathUnescape(cr.Name)
		filePath := os.Getenv("STATIC_FILES_PATH") + "/" + relpath

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

func getTorrentFromHash(service torrenter.UseCase, hash string) (*torrent.Torrent, bool) {
	var hexHash [20]byte

	hex, err := hex.DecodeString(hash)
	if err != nil {
		return nil, false
	}
	copy(hexHash[:], hex[:20])
	return service.Torrent(hexHash)
}

func HLSHandler(service iostream.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("filename")
		dirName := c.Param("dirname")

		err := service.Save(dirName, fileName, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}
