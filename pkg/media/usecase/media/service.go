package media

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"log"
	"math"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
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

func (s *Service) StreamMagnet(magnet string) (*entity.StreamResponse, error) {
	tc := s.tc

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}

	<-t.GotInfo()

	info := t.Info()
	if info.Files != nil {
		return nil, errors.New("multiple file torrent")
	}

	// set higher priority to the first 1% pieces
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
		return nil, err
	}

	c.WaitUntilReady()

	go func() {
		c.WaitUntilDone()
		err := c.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	url := "http://localhost" + ":" + os.Getenv("MEDIA_PORT") + "/" + filepath
	rv := entity.NewStreamResponse(t.Name(), t.InfoHash().HexString(), url)
	return rv, nil
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
