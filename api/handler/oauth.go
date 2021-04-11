package handler

import (
	"net/http"
	"time"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/config"
	auth "github.com/MonsieurTa/hypertube/usecase/authentication"
	"github.com/MonsieurTa/hypertube/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/usecase/state"
	"github.com/gin-gonic/gin"
)

func MakeOAuth2Handlers(
	g *gin.RouterGroup,
	ftService fortytwo.UseCase,
	stateService state.UseCase,
	service auth.UseCase,
) {
	g.POST("/token", accessTokenGeneration(service))

	g.GET("/fortytwo/callback", redirectCallback(ftService, stateService))
	g.GET("/fortytwo/authorize_uri", getAuthorizeURI(ftService, stateService))
}

func accessTokenGeneration(service auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUserCredentialValidator()

		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewValidationError(err))
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
