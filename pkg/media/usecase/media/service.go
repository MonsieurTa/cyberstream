package media

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/MonsieurTa/hypertube/pkg/media/internal/hls"
	"github.com/anacrolix/torrent"
)

type Service struct {
	tc *torrent.Client
}

func NewService(tc *torrent.Client) *Service {
	return &Service{
		tc,
	}
}

func (s *Service) StreamMagnet(magnet string) (string, <-chan bool, error) {
	tc := s.tc

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		return "", nil, err
	}

	<-t.GotInfo()

	h := sha1.Sum([]byte(t.Info().Name))
	dir := hex.EncodeToString(h[:])
	filepath := dir + "/" + "out.m3u8"
	hlsPath := os.Getenv("STATIC_FILES_PATH") + "/" + filepath
	fmt.Println(hlsPath)

	os.Mkdir(os.Getenv("STATIC_FILES_PATH")+"/"+dir, os.ModePerm)

	ready, done := toHLS(t, hlsPath)
	go func() {
		<-done
		close(ready)
		close(done)
	}()

	return filepath, ready, nil
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
