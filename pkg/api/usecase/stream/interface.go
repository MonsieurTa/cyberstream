package stream

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

type Reader interface {
	FindByID(movieID uuid.UUID) (*entity.Movie, error)
}

type Writer interface {
	Create(movie *entity.Movie) (uuid.UUID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Stream(m *entity.Movie) (string, error)
}
