package jackett

import "github.com/MonsieurTa/hypertube/common/infrastructure/repository"

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{repo}
}

func (s *Service) AllIndexers() (*repository.Indexers, error) {
	return s.repo.Indexers(false)
}

func (s *Service) ConfiguredIndexers() (*repository.Indexers, error) {
	return s.repo.Indexers(true)
}

func (s *Service) Search() {
}
