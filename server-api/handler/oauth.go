package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/server-api/common"
	auth "github.com/MonsieurTa/hypertube/server-api/usecase/authentication"
	"github.com/gin-gonic/gin"
)

func AccessTokenGeneration(service auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUserCredentialsValidator()

		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("validation", err))
			return
		}

		input := validator.Value()

		userID, err := service.Authenticate(input.Username, input.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, common.NewError("auth", err))
			return
		}

		expiresAt := time.Now().Add(5 * time.Minute)
		token := service.GenerateAccessToken(userID, expiresAt)
		secret := os.Getenv("JWT_SECRET")
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewError("auth", err))
		}
		c.Header("Authorization", `Bearer `+tokenString)
	}
}
