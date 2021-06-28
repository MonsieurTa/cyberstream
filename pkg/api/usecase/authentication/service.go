package authentication

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	errInvalidClaims    = errors.New("invalid claims")
	errInvalidTokenType = errors.New("invalid tokenType")
	errInvalidToken     = errors.New("invalid token")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{
		repo,
	}
}

func (s *Service) Authenticate(username, password string) (uuid.UUID, error) {
	return s.repo.CredentialsExist(username, password)
}

func (s *Service) NewRefreshToken(userID uuid.UUID) (string, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	return s.newToken("refresh", userID, exp)
}

func (s *Service) NewAccessToken(userID uuid.UUID) (string, error) {
	exp := time.Now().Add(15 * time.Minute).Unix()
	return s.newToken("access", userID, exp)
}

func (s *Service) newToken(tokenType string, userID uuid.UUID, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"type":    tokenType,
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Service) ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_TOKEN")
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		log.Println(err.Error())
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errInvalidClaims
	}

	tokenType, ok := claims["type"]
	if !ok || tokenType.(string) != "refresh" {
		return nil, errInvalidTokenType
	}
	return token, nil
}

func (s *Service) ExtractMetadata(token *jwt.Token, tokenType string) (uuid.UUID, error) {
	if err := validMethod(token); err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenType, ok := claims["type"]
		if !ok {
			return uuid.Nil, errInvalidClaims
		}
		if tokenType != tokenType {
			return uuid.Nil, errInvalidTokenType
		}

		userID, ok := claims["user_id"]
		if !ok {
			return uuid.Nil, errInvalidClaims
		}
		return uuid.Parse(userID.(string))
	}
	return uuid.Nil, errInvalidToken
}

func (s *Service) UserExists(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	return err
}

func validMethod(token *jwt.Token) error {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return nil
}
