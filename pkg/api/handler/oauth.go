package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	auth "github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/dgrijalva/jwt-go"
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

		t, err := service.NewToken(userID.String(), "")
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusBadRequest, cerr)
			return
		}
		c.SetCookie("refresh_token", t.RefreshToken, 60*60*24, "/", "", true, true)
		c.Header("Authorization", "Bearer "+t.AccessToken)
		c.Status(http.StatusOK)
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
