package jackett

import "github.com/MonsieurTa/hypertube/common/infrastructure/repository"

type Reader interface {
	Indexers(configured bool) (*repository.Indexers, error)
}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	AllIndexers() (*repository.Indexers, error)
	ConfiguredIndexers() (*repository.Indexers, error)
	Search()
}
