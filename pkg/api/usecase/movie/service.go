package movie

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{repo}
}

func (s Service) FindByID(id uuid.UUID) (*entity.Movie, error) {
	return s.repo.FindByID(id)
}
func (s Service) Register(movie *entity.Movie) (uuid.UUID, error) {
	return s.repo.Create(movie)
}
