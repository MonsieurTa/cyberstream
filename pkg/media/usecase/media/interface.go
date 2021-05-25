package media

import "github.com/MonsieurTa/hypertube/common/entity"

type Reader interface{}
type Writer interface {
	AddMagnet()
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	StreamMagnet(magnet string) (*entity.StreamResponse, error)
}
