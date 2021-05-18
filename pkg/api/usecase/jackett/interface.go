package jackett

import (
	"context"

	"github.com/webtor-io/go-jackett"
)

type Reader interface {
	Fetch(ctx context.Context, fr *jackett.FetchRequest) (*jackett.FetchResponse, error)
}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Search(pattern string, categories []uint) (*jackett.FetchResponse, error)
	Categories() map[string]uint
}
