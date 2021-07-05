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
	NewToken(id, from string) (*Token, error)
	NewRefreshToken(id, from string) (string, error)
	NewAccessToken(id, from string) (string, error)
	ExtractMetadata(claims jwt.MapClaims, tokenType string) (TokenMeta, error)
	UserExists(id uuid.UUID) error
	ValidateRefreshToken(tokenStr string) (*jwt.Token, error)
	// ValidRefreshToken(token string) bool
	// ValidAccessToken(token string) bool
}
