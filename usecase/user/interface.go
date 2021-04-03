package user

import (
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/google/uuid"
)

type Reader interface{}

type Writer interface {
	Create(e *entity.User) (uuid.UUID, error)
}

type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	RegisterUser(c entity.User) (uuid.UUID, error)
}
