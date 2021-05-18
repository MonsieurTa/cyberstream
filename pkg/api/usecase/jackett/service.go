package jackett

import (
	"context"

	"github.com/webtor-io/go-jackett"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{repo}
}

func (s *Service) fetch(fr *jackett.FetchRequest) (*jackett.FetchResponse, error) {
	ctx := context.Background()
	return s.repo.Fetch(ctx, fr)
}

func (s *Service) Categories() map[string]uint {
	return map[string]uint{
		"TV":     TV,
		"Anime":  TV_ANIME,
		"Movies": MOVIES,
	}
}

func (s *Service) Search(pattern string, categories []uint) (*jackett.FetchResponse, error) {
	fr := &jackett.FetchRequest{
		Query:              pattern,
		Categories:         categories,
		ConfiguredIndexers: true,
	}
	return s.fetch(fr)
}
