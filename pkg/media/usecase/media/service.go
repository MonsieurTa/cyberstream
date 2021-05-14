package media

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"math"
	"os"

	"github.com/MonsieurTa/hypertube/pkg/media/internal/hls"
	"github.com/anacrolix/torrent"
)

type Service struct {
	tc *torrent.Client
}

func NewService(tc *torrent.Client) UseCase {
	return &Service{
		tc,
	}
}

func (s *Service) StreamMagnet(magnet string) (string, error) {
	tc := s.tc

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		return "", err
	}

	<-t.GotInfo()

	// set higher priority to the first 1% pieces
	info := t.Info()
	numPieces := info.NumPieces()
	threshold := int(math.Ceil(float64(numPieces) / 100))
	for i := 0; i < threshold; i++ {
		t.Piece(i).SetPriority(torrent.PiecePriorityNow)
	}

	dirName, filepath, hlspath := stringify(t.Info().Name)

	createHLSFolder(dirName)

	c := hls.NewHLSConverter(hlspath, t)
	err = c.Convert()
	if err != nil {
		return "", err
	}

	c.WaitUntilReady()

	go func() {
		c.WaitUntilDone()
		err := c.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	return filepath, nil
}

func stringify(name string) (dirName, filepath, hlspath string) {
	h := sha1.Sum([]byte(name))
	dirName = hex.EncodeToString(h[:])
	filepath = dirName + "/" + "out.m3u8"
	hlspath = os.Getenv("STATIC_FILES_PATH") + "/" + filepath
	return
}

func createHLSFolder(dirName string) {
	os.Mkdir(os.Getenv("STATIC_FILES_PATH")+"/"+dirName, os.ModePerm)
}
