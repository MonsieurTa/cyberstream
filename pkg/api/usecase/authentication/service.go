package authentication

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	errInvalidClaims    = errors.New("invalid claims")
	errInvalidTokenType = errors.New("invalid tokenType")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{
		repo,
	}
}

type TokenMeta map[string]string

func (s *Service) Authenticate(username, password string) (uuid.UUID, error) {
	return s.repo.CredentialsExist(username, password)
}

func (s *Service) ExtractMetadata(claims jwt.MapClaims, expectedType string) (TokenMeta, error) {
	rv := make(TokenMeta)
	rv["from"] = ""

	tokenType, ok := claims["type"]
	if !ok {
		return nil, errInvalidClaims
	}
	if tokenType != expectedType {
		return nil, errInvalidTokenType
	}

	ID, ok := claims["user_id"]
	if !ok {
		return nil, errInvalidClaims
	}

	from, ok := claims["from"]
	if ok {
		rv["from"] = from.(string)
	}

	rv["type"] = tokenType.(string)
	rv["user_id"] = ID.(string)
	return rv, nil
}

func (s *Service) UserExists(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	return err
}
