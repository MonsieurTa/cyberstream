package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Reader interface {
	Validate(username, password string) error
}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Authenticate(username, password string) error
	GenerateAccessToken(username string, expiresAt time.Time) *jwt.Token
}
