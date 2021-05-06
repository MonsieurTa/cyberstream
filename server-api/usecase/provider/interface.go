package provider

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

type Reader interface {
	FindByName(name entity.ProviderName) (*entity.Provider, error)
}

type Writer interface {
	RegisterProviders(providers []entity.Provider) error
	StoreMovie(provider *entity.Provider, movie *entity.Movie) (uuid.UUID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	StoreMovie(name entity.ProviderName, movie *entity.Movie) (uuid.UUID, error)
}
