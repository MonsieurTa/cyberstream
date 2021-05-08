package validator

import (
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/gin-gonic/gin"
)

type UserCredentialsValidator struct {
	Credentials struct {
		Username string `json:"username" binding:"required,alphanum,min=4,max=255"`
		Password string `json:"password" binding:"required,alphanum,min=8,max=255"`
	} `json:"credentials"`

	output CredentialsOutput `json:"-"`
}

type CredentialsOutput struct {
	Username string
	Password string
}

func NewUserCredentialsValidator() *UserCredentialsValidator {
	return &UserCredentialsValidator{}
}

func (v *UserCredentialsValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	v.output.Username = v.Credentials.Username
	v.output.Password = v.Credentials.Password
	return nil
}

func (v UserCredentialsValidator) Value() CredentialsOutput {
	return v.output
}
