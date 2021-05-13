package subsplease

import (
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
)

type Service struct {
	repo Repository

	endpoint string
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		endpoint: os.Getenv("SUBSPLEASE_API_URL"),
	}
}

func (s *Service) Latests() ([]entity.Movie, error) {
	// subsPleaseEpisodes, err := s.repo.Latest()
	_ , err := s.repo.Latest()
	if err != nil {
		return []entity.Movie{}, err
	}
	// TODO encrypt magnets and replace them in returned values (i.e AES Cypher)
	return []entity.Movie{}, nil
}
