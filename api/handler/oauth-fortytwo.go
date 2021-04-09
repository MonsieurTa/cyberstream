package handler

import (
	"errors"
	"net/http"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/infrastructure/repository"
	"github.com/MonsieurTa/hypertube/usecase/fortytwo"
	"github.com/MonsieurTa/hypertube/usecase/state"
	"github.com/gin-gonic/gin"
)

func redirectCallback(ftService fortytwo.UseCase, stateService state.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: create a specific validator
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			c.JSON(http.StatusBadRequest, common.NewError("auth", errors.New("invalid_parameter")))
			return
		}

		if err := stateService.Validate(state); err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusUnauthorized, cerr)
			return
		}
		stateService.Delete(state)

		token, err := ftService.GetAccessToken(code, state)
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusUnauthorized, cerr)
			return
		}
		c.JSON(http.StatusOK, token)
	}
}

func getAuthorizeURI(ftService fortytwo.UseCase, stateService state.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := repository.GenerateState()
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusInternalServerError, cerr)
			return
		}

		uri, err := ftService.GetAuthorizeURI(state)
		if err != nil {
			cerr := common.NewError("auth", err)
			c.JSON(http.StatusBadRequest, cerr)
			return
		}
		stateService.Save(state)
		c.JSON(http.StatusOK, uri)
	}
}
