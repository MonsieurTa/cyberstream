package validator

import (
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/gin-gonic/gin"
)

type SearchValidator struct {
	SearchParams struct {
		Pattern    string `form:"pattern" json:"pattern" binding:"required"`
		Categories []uint `form:"category" json:"category" binding:"required"`
	}

	Pattern    string
	Categories []uint
}

func NewSearchValidator() *SearchValidator {
	return &SearchValidator{}
}

func (v *SearchValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	v.Pattern = v.SearchParams.Pattern
	v.Categories = v.SearchParams.Categories
	return nil
}

func (v *SearchValidator) Value() (string, []uint) {
	return v.Pattern, v.Categories
}
