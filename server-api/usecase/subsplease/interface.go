package subsplease

import "github.com/MonsieurTa/hypertube/common/entity"

type Reader interface {
}

type Repository interface {
	Reader
}

type UseCase interface {
	Latests() ([]entity.Movie, error)
}
