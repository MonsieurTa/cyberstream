package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/media/internal/hls"
	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
)

const (
	ERR_MAGNET         = "error_magnet"
	ERR_TORRENT_CLIENT = "error_torrent_client"
	ERR_VALIDATION     = "error_validation"
)

// TODO create Stream service
func Stream(c *gin.Context) {
	v := validator.NewStreamRequestValidator()

	err := v.Validate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_VALIDATION: err.Error(),
		})
		return
	}

	magnet := v.Value().Magnet

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = os.Getenv("DOWNLOAD_FILES_PATH")
	tc, err := torrent.NewClient(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			ERR_TORRENT_CLIENT: err.Error(),
		})
		return
	}

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_MAGNET: err.Error(),
		})
		tc.Close()
		return
	}

	<-t.GotInfo()

	h := sha1.Sum([]byte(t.Info().Name))
	dir := hex.EncodeToString(h[:])
	filepath := dir + "/" + "out.m3u8"
	hlsPath := os.Getenv("STATIC_FILES_PATH") + "/" + filepath

	os.Mkdir(os.Getenv("STATIC_FILES_PATH")+"/"+dir, os.ModePerm)

	ready, done := toHLS(t, hlsPath)
	go func() {
		<-done
		close(ready)
		close(done)
		tc.Close()
	}()

	<-ready

	resp := entity.StreamResponse{
		Url: "http://localhost" + ":" + os.Getenv("MEDIA_PORT") + "/" + filepath,
	}
	c.JSON(http.StatusOK, &resp)
}

func toHLS(t *torrent.Torrent, path string) (chan bool, chan bool) {
	rpipe, wpipe, wait := hls.Init(path)
	ready := make(chan bool)
	done := make(chan bool)

	progress := make(chan int64)
	go func() {
		r := t.NewReader()
		buf := make([]byte, t.Info().PieceLength)
		at := int64(0)
		end := t.Length()

		r.SetReadahead(end / 100 * 5)
		for at < end {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			_, err = wpipe.Write(buf)
			if err != nil {
				log.Fatal(err)
			}
			at += int64(n)
			select {
			case progress <- at:
			default:
			}
		}
		done <- true
		log.Println(path, ": done")

		r.Close()
		rpipe.Close()
		wpipe.Close()
		close(progress)
		wait()
	}()

	// send true to ready when 1% downloaded
	go func() {
		threshold := t.Length() / 100
		for at := range progress {
			if at >= threshold {
				break
			}
			time.Sleep(time.Microsecond * 100)
		}
		ready <- true
		log.Println(path, ": ready")
	}()
	return ready, done
}
