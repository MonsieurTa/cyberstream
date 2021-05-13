package subsplease

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/pkg/api/internal/subsplease"
)

type Reader interface {
	Latest() ([]subsplease.Episode, error)
}

type Repository interface {
	Reader
}

type UseCase interface {
	Latest() ([]entity.Movie, error)
}
