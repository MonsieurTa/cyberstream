package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service struct {
	repo Repository
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewService(repo Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) Authenticate(username, password string) error {
	return s.repo.Validate(username, password)
}

func (s *Service) GenerateAccessToken(username string, expiresAt time.Time) *jwt.Token {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
