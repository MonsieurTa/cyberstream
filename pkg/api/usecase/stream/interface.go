package stream

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

type Reader interface {
	FindByID(videoID uuid.UUID) (*entity.Video, error)
	FindByName(name string) (*entity.Video, error)
	FindByHash(hash string) (*entity.Video, error)
}

type Writer interface {
	Create(video *entity.Video) (uuid.UUID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Stream(streamReq *entity.StreamRequest) (*entity.StreamResponse, error)
}
