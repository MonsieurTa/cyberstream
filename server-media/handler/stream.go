package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/internal/hls"
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
	cfg.DataDir = "/download"
	tc, err := torrent.NewClient(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			ERR_TORRENT_CLIENT: err.Error(),
		})
		return
	}
	defer tc.Close()

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			ERR_MAGNET: err.Error(),
		})
		return
	}

	<-t.GotInfo()

	h := sha1.Sum([]byte(t.Info().Name))
	path := hex.EncodeToString(h[:]) + `.m3u8`

	ready := toHLS(t, path)
	defer close(ready)

	fmt.Println("====== NOT READY YET ======")
	<-ready
	fmt.Println("====== READY SENDING RESPONSE ======")

	resp := entity.StreamResponse{
		Url: "localhost:3001" + config.STATIC_FILES_PATH + "/" + path,
	}
	c.JSON(http.StatusOK, &resp)
}

func toHLS(t *torrent.Torrent, path string) chan bool {
	rpipe, wpipe, wait := hls.Init(path)

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
		r.Close()
		close(progress)
		rpipe.Close()
		wpipe.Close()
		wait()
	}()

	ready := make(chan bool)

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
	}()
	return ready
}
