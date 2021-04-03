package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/api/validator"
	"github.com/MonsieurTa/hypertube/usecase/user"
	"github.com/gin-gonic/gin"
)

func MakeUsersHandlers(g *gin.RouterGroup, service user.UseCase) {
	g.POST("/", usersRegistration(service))
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
		ID, err := service.CreateUser(input)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
		c.JSON(http.StatusOK, ID)
	}
}
