package validator

import (
	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/gin-gonic/gin"
)

type UserRegistrationValidator struct {
	User struct {
		PublicInfo struct {
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
			Phone     string `json:"phone"`
			Email     string `json:"email"`
		} `json:"public_info"`
		Credential struct {
			Username string `json:"Username"`
			Password string `json:"password"`
		} `json:"credential"`
	} `json:"user"`

	output *entity.User `json:"-"`
}

func NewUserRegistrationValidator() *UserRegistrationValidator {
	return &UserRegistrationValidator{}
}

func (v *UserRegistrationValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}
	input := entity.CreateUserT{
		FirstName: v.User.PublicInfo.FirstName,
		LastName:  v.User.PublicInfo.LastName,
		Phone:     v.User.PublicInfo.Phone,
		Email:     v.User.PublicInfo.Email,
		Username:  v.User.Credential.Username,
		Password:  v.User.Credential.Password,
	}
	v.output, err = entity.NewUser(input)
	if err != nil {
		return err
	}
	return nil
}

func (v UserRegistrationValidator) Value() entity.User {
	return *v.output
}
