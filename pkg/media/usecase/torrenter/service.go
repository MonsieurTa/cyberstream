package torrenter

import (
	"log"

	"github.com/MonsieurTa/hypertube/common/entity"
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

func (s *Service) DownloadMagnet(streamReq *entity.StreamRequest) (*streaminfo.StreamInfo, error) {
	tc := s.tc

	t, err := tc.AddMagnet(streamReq.Magnet)
	if err != nil {
		return nil, err
	}

	log.Println("Getting torrent info...")
	<-t.GotInfo()

	streamInfo, err := streaminfo.Extract(t)
	if err != nil {
		return nil, err
	}
	streamInfo.InfoHash = streamReq.InfoHash

	if streamInfo.HasSubtitles() {
		for _, subf := range streamInfo.SubtitlesFiles {
			subf.SetPriority(torrent.PiecePriorityNow)
		}
	}

	streamInfo.StreamFile.SetPriority(torrent.PiecePriorityReadahead)
	return streamInfo, nil
}
