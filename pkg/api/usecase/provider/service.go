package provider

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

var PROVIDERS = []entity.Provider{
	{Name: entity.SUBSPLEASE},
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) (*Service, error) {
	err := repo.RegisterProviders(PROVIDERS)
	if err != nil {
		return nil, err
	}
	return &Service{repo}, nil
}

func (s *Service) StoreVideo(name entity.ProviderName, video *entity.Video) (uuid.UUID, error) {
	provider, err := s.repo.FindByName(name)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := s.repo.StoreVideo(provider, video)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
