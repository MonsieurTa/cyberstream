package user

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
)

type Reader interface {
	CredentialsExist(username, password string) (uuid.UUID, error)
}

type Writer interface {
	Create(e *entity.User) (uuid.UUID, error)
	UpdateCredentials(userID uuid.UUID, username, password string) error
	UpdatePublicInfo(userID uuid.UUID, email, pictureURL string) error
}

type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	Register(c entity.User) (uuid.UUID, error)
	UpdateCredentials(userID uuid.UUID, username, password string) error
	UpdatePublicInfo(userID uuid.UUID, email, pictureURL string) error
}
