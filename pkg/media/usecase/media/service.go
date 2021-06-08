package media

import (
	"log"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/pkg/media/internal/hls"
	"github.com/MonsieurTa/hypertube/pkg/media/internal/streaminfo"
	"github.com/anacrolix/torrent"
)

type Service struct {
	tc *torrent.Client
}

func NewService(tc *torrent.Client) UseCase {
	return &Service{tc}
}

func (s *Service) Torrent(hash [20]byte) (*torrent.Torrent, bool) {
	return s.tc.Torrent(hash)
}

func (s *Service) StreamMagnet(magnet string) *entity.StreamResponse {
	var resp entity.StreamResponse
	baseURL := "http://localhost" + ":" + os.Getenv("MEDIA_PORT")
	tc := s.tc

	t, err := tc.AddMagnet(magnet)
	if err != nil {
		resp.Error = err.Error()
		return &resp
	}

	log.Println("Getting torrent info...")
	<-t.GotInfo()

	streamInfo, err := streaminfo.Extract(t)
	if err != nil {
		resp.Error = err.Error()
		return &resp
	}
	resp.Name = streamInfo.StreamFile.DisplayPath()
	resp.InfoHash = streamInfo.InfoHash
	resp.Ext = streamInfo.StreamFile.Ext()

	if streamInfo.HasSubtitles() {
		resp.SubtitlesURLs = make([]string, len(streamInfo.SubtitlesFiles))
		for i, subf := range streamInfo.SubtitlesFiles {
			subf.SetPriority(torrent.PiecePriorityNow)
			resp.SubtitlesURLs[i] = baseURL + "/static/" + subf.Path()
		}
	}

	if streamInfo.StreamFile.Ext() == ".mp4" {
		streamInfo.StreamFile.SetPriority(torrent.PiecePriorityReadahead)
		resp.MediaURL = baseURL + "/content?name=" + streamInfo.StreamFile.Path()
		return &resp
	}

	log.Println("Starting conversion...")
	playlistPath, err := convertToHLS(streamInfo.InfoHash, streamInfo.StreamFile.File)
	if err != nil {
		resp.Error = err.Error()
		return &resp
	}
	resp.MediaURL = baseURL + "/static/" + playlistPath
	return &resp
}

func convertToHLS(infoHash string, f *torrent.File) (string, error) {
	c := hls.NewHLSConverter(&hls.Config{
		StreamFile: f,
		OutputDir:  os.Getenv("STATIC_FILES_PATH") + "/" + infoHash,
	})

	c.Convert()

	go func() {
		c.WaitUntilDone()
		c.Close()
	}()

	err := c.WaitUntilReady()
	if err != nil {
		return "", err
	}

	return infoHash + "/master.m3u8", nil
}
