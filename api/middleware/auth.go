package middleware

import (
	"net/http"

	auth "github.com/MonsieurTa/hypertube/usecase/authentication"
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

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &auth.Claims{}

		_, err := request.ParseFromRequestWithClaims(c.Request, request.OAuth2Extractor, claims, secretGiver(secret))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
	}
}
