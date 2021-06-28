package authentication

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Reader interface {
	CredentialsExist(username, password string) (uuid.UUID, error)
	FindByID(userID uuid.UUID) (*entity.User, error)
}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Authenticate(username, password string) (uuid.UUID, error)
	NewRefreshToken(userID uuid.UUID) (string, error)
	NewAccessToken(userID uuid.UUID) (string, error)
	ExtractMetadata(token *jwt.Token, tokenType string) (uuid.UUID, error)
	UserExists(id uuid.UUID) error
	ValidateRefreshToken(tokenStr string) (*jwt.Token, error)
	// ValidRefreshToken(token string) bool
	// ValidAccessToken(token string) bool
}
