package user

import (
	"github.com/MonsieurTa/hypertube/common/entity"
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

func (s *Service) Register(model entity.User) (uuid.UUID, error) {
	ID, err := s.repo.Create(&model)
	if err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}

func (s *Service) UpdateCredentials(userID uuid.UUID, username, password string) error {
	return s.repo.UpdateCredentials(userID, username, password)
}

func (s *Service) UpdatePublicInfo(userID uuid.UUID, email, pictureURL string) error {
	return s.repo.UpdatePublicInfo(userID, email, pictureURL)
}
