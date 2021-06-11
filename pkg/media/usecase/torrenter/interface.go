package torrenter

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/pkg/media/internal/streaminfo"
	"github.com/anacrolix/torrent"
)

type Reader interface{}
type Writer interface {
	AddMagnet()
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	DownloadMagnet(streamReq *entity.StreamRequest) (*streaminfo.StreamInfo, error)
	Torrent(hash [20]byte) (*torrent.Torrent, bool)
}
