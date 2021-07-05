package authentication

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) NewToken(id, from string) (*Token, error) {
	accessExp := time.Now().Add(15 * time.Minute).Unix()
	refreshExp := time.Now().Add(24 * time.Hour).Unix()
	at, err := s.newToken("access", id, from, accessExp)
	if err != nil {
		return nil, err
	}
	rt, err := s.newToken("refresh", id, from, refreshExp)
	if err != nil {
		return nil, err
	}
	return &Token{at, rt}, nil
}

func (s *Service) NewRefreshToken(id, from string) (string, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	return s.newToken("refresh", id, from, exp)
}

func (s *Service) NewAccessToken(id, from string) (string, error) {
	exp := time.Now().Add(15 * time.Minute).Unix()
	return s.newToken("access", id, from, exp)
}

func (s *Service) newToken(tokenType, id, from string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"type":    tokenType,
		"user_id": id,
		"from":    from,
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
		secret := os.Getenv("JWT_SECRET")
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
