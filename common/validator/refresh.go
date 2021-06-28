package validator

import (
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/gin-gonic/gin"
)

type RefreshValidator struct {
	Refresh struct {
		RefreshToken string `json:"refresh_token" binding:"required,alphanum,min=36,max=36"`
	}
	output string `json:"-"`
}

func NewRefreshValidator() *RefreshValidator {
	return &RefreshValidator{}
}

func (v *RefreshValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	v.output = v.Refresh.RefreshToken
	return nil
}

func (v RefreshValidator) Value() string {
	return v.output
}
