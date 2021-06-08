package media

import (
	"github.com/MonsieurTa/hypertube/common/entity"
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
	StreamMagnet(magnet string) *entity.StreamResponse
	Torrent(hash [20]byte) (*torrent.Torrent, bool)
}
