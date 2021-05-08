package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) Authenticate(username, password string) (uuid.UUID, error) {
	return s.repo.CredentialsExist(username, password)
}

func (s *Service) GenerateAccessToken(userID uuid.UUID, expiresAt time.Time) *jwt.Token {
	claims := jwt.StandardClaims{
		Audience:  userID.String(),
		ExpiresAt: expiresAt.Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
