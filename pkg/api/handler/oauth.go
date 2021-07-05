package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	auth "github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/gin-gonic/gin"
)

func Login(service auth.UseCase) gin.HandlerFunc {
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

		refreshToken, err := service.NewRefreshToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewError("auth", err))
			return
		}
		accessToken, err := service.NewAccessToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewError("auth", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"refresh_token": refreshToken,
			"access_token":  accessToken,
		})
	}
}

func AccessTokenGeneration(service auth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtStr, err := c.Cookie("refresh_token")
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		rt, err := service.ValidateRefreshToken(jwtStr)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		meta, err := service.ExtractMetadata(rt.Claims.(jwt.MapClaims), "refresh")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		accessToken, err := service.NewAccessToken(meta["id"], meta["from"])
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Header("Authorization", "Bearer "+accessToken)
		c.Status(http.StatusOK)
	}
}
