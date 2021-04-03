package user

import (
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) RegisterUser(model entity.User) (uuid.UUID, error) {
	id, err := s.repo.Create(&model)
	if err != nil {
		return id, err
	}
	return id, nil
}
