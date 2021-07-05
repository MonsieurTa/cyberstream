package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func secretGiver(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(secret))
		return b, nil
	}
}

func Auth(service authentication.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authField, ok := c.Request.Header["Authorization"]

		if !ok || len(authField) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s := strings.Split(authField[0], "Bearer ")
		if len(s) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err := validateToken(s[1], service)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func validateToken(tokenStr string, service authentication.UseCase) error {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, secretGiver(secret))
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	meta, err := service.ExtractMetadata(claims, "access")
	if err != nil {
		return err
	}

	from, ok := meta["from"]
	if ok {
		if from != "42" && from != "" {
			return err
		}
		return nil
	}

	id, err := uuid.Parse(meta["user_id"])
	if err != nil {
		return err
	}

	err = service.UserExists(id)
	if err != nil {
		return err
	}
	return nil
}
