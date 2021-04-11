package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/api/middleware"
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/config"
	"github.com/MonsieurTa/hypertube/usecase/user"
	"github.com/gin-gonic/gin"
)

func MakeUsersHandlers(g *gin.RouterGroup, service user.UseCase) {
	g.POST("/", usersRegistration(service))
	g.PATCH("/", middleware.Auth(config.JWT_SECRET), usersUpdate(service))
}

func usersRegistration(service user.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUserRegistrationValidator()

		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewValidationError(err))
			return
		}

		input := validator.Value()
		ID, err := service.RegisterUser(input)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
		c.JSON(http.StatusOK, ID)
	}
}

func usersUpdate(service user.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewUserUpdateValidator()

		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("validation", err))
			return
		}

		input := validator.Value()
		err = service.UpdateCredentials(&input.UserID, input.Username, input.Password)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
		err = service.UpdatePublicInfo(&input.UserID, input.Email, input.PictureURL)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}
