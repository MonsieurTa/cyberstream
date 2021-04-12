package validator

import (
	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/gin-gonic/gin"
)

type UserCredentialValidator struct {
	Credentials struct {
		Username string `json:"username" binding:"required,alphanum,min=4,max=255"`
		Password string `json:"password" binding:"required,alphanum,min=8,max=255"`
	} `json:"credentials"`

	output CredentialOutput `json:"-"`
}

type CredentialOutput struct {
	Username string
	Password string
}

func NewUserCredentialValidator() *UserCredentialValidator {
	return &UserCredentialValidator{}
}

func (v *UserCredentialValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	v.output.Username = v.Credentials.Username
	v.output.Password = v.Credentials.Password
	return nil
}

func (v UserCredentialValidator) Value() CredentialOutput {
	return v.output
}
