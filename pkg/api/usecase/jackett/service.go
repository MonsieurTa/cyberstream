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

func (s *Service) Fetch(fr *jackett.FetchRequest) (*jackett.FetchResponse, error) {
	ctx := context.Background()
	return s.repo.Fetch(ctx, fr)
}
