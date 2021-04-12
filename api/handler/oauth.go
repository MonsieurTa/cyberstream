package handler

import (
	"net/http"
	"time"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/config"
	auth "github.com/MonsieurTa/hypertube/usecase/authentication"
	"github.com/gin-gonic/gin"
)

func AccessTokenGeneration(service auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUserCredentialValidator()

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
		tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewError("auth", err))
		}
		c.Header("Authorization", `Bearer `+tokenString)
	}
}
