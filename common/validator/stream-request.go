package validator

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/server-api/common"
	"github.com/gin-gonic/gin"
)

type StreamRequestValidator struct {
	StreamRequest struct {
		Magnet string `json:"magnet" bindind:"required,min=1"`
	} `json:"stream_request"`

	output entity.StreamRequest `json:"-"`
}

func NewStreamRequestValidator() *StreamRequestValidator {
	return &StreamRequestValidator{}
}

func (v *StreamRequestValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}
	v.output.Magnet = v.StreamRequest.Magnet
	return nil
}

func (v StreamRequestValidator) Value() entity.StreamRequest {
	return v.output
}
