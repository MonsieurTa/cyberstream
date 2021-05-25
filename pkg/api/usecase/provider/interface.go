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
	StoreVideo(provider *entity.Provider, video *entity.Video) (uuid.UUID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	StoreVideo(name entity.ProviderName, video *entity.Video) (uuid.UUID, error)
}
