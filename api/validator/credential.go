package validator

import (
	"github.com/MonsieurTa/hypertube/api/common"
	"github.com/gin-gonic/gin"
)

type UserCredentialValidator struct {
	Credential struct {
		Username string `json:"username"`
		Password string `json:"password"`
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

	v.output.Username = v.Credential.Username
	v.output.Password = v.Credential.Password
	return nil
}

func (v UserCredentialValidator) Value() CredentialOutput {
	return v.output
}
