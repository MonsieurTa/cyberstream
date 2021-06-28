package middleware

import (
	"errors"
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func secretGiver(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(secret))
		return b, nil
	}
}

func Auth(service authentication.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")
		claims := &jwt.MapClaims{}

		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, secretGiver(secret), request.WithClaims(claims))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
			return
		}

		userID, err := service.ExtractMetadata(token, "access")
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		err = service.UserExists(userID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}
}
