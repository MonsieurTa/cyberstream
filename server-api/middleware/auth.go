package middleware

import (
	"net/http"

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
		claims := &jwt.StandardClaims{}

		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, secretGiver(secret), request.WithClaims(claims))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		c.Set("token", claims)
	}
}
