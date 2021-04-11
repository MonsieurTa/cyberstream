package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Reader interface {
	CredentialsExist(username, password string) (*uuid.UUID, error)
}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Authenticate(username, password string) (*uuid.UUID, error)
	GenerateAccessToken(userID *uuid.UUID, expiresAt time.Time) *jwt.Token
}
