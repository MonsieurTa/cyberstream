package handler

import (
	"errors"
	"net/http"

	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/authentication"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/fortytwo"
	"github.com/gin-gonic/gin"
)

func RedirectCallback(ft fortytwo.UseCase, auth authentication.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: create a specific validator
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			c.JSON(http.StatusBadRequest, common.NewError("auth", errors.New("invalid_parameter")))
			return
		}

		token, err := ft.GetToken(code, state)
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusUnauthorized, cerr)
			return
		}

		userInfo, err := ft.GetUserInfo(token.AccessToken)
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusBadRequest, cerr)
			return
		}

		t, err := auth.NewToken(userInfo.Login, "42")
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

func GetAuthorizeURI(ftService fortytwo.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		uri, err := ftService.GetAuthorizeURI()
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusBadRequest, cerr)
			return
		}
		c.JSON(http.StatusOK, uri)
	}
}
