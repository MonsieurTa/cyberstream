package handler

import (
	"errors"
	"net/http"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/usecase/fortytwo"
	"github.com/gin-gonic/gin"
)

func MakeFortyTwoAuthHandlers(g *gin.RouterGroup, service fortytwo.UseCase) {
	g.GET("/fortytwo/callback", redirectCallback(service))
	g.GET("/fortytwo/authorize_uri", getAuthorizeURI(service))
}

func redirectCallback(service fortytwo.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, errors.New("code invalid"))
			return
		}
		state := c.Query("state")
		if state == "" {
			c.JSON(http.StatusBadRequest, errors.New("state invalid"))
			return
		}

		token, err := service.GetAccessToken(code, state)
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusUnauthorized, cerr)
			return
		}
		c.JSON(http.StatusOK, token)
	}
}

func getAuthorizeURI(service fortytwo.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		uri, err := service.GetAuthorizeURI()
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusBadRequest, cerr)
			return
		}
		c.JSON(http.StatusOK, uri)
	}
}
