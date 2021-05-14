package subsplease

import (
	"log"
	"os"

	"github.com/MonsieurTa/hypertube/common/cipher"
	"github.com/MonsieurTa/hypertube/common/entity"
)

type Service struct {
	repo Repository

	endpoint string
}

func NewService(repo Repository) UseCase {
	return &Service{
		repo:     repo,
		endpoint: os.Getenv("SUBSPLEASE_API_URL"),
	}
}

func (s *Service) Latest() ([]*entity.Movie, error) {
	episodes, err := s.repo.Latest()
	if err != nil {
		return []*entity.Movie{}, err
	}

	t, err := cipher.NewCryptograph(os.Getenv("AES_KEY"))
	if err != nil {
		return []*entity.Movie{}, err
	}

	rv := make([]*entity.Movie, 0, len(episodes))
	for _, e := range episodes {
		plainMagnet := []byte(e.HighestResolutionMagnet())
		encrypted, err := t.Encrypt(plainMagnet)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		name := e.Show + " - Episode " + e.Episode
		rv = append(rv, entity.NewMovie(name, "", encrypted))
	}
	return rv, nil
}
