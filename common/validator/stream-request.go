package validator

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/MonsieurTa/hypertube/common/cipher"
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/pkg/api/common"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidMagnet = errors.New("invalid magnet")
)

type StreamRequestValidator struct {
	StreamRequest struct {
		Name   string `form:"name" json:"name" bindind:"required,min=1"`
		Magnet string `form:"magnet" json:"magnet" bindind:"required,min=1"`
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

	decryptedMagnet, err := decryptMagnet(v.StreamRequest.Magnet)
	if err != nil {
		return err
	}

	v.output.Name = v.StreamRequest.Name
	v.output.Magnet = decryptedMagnet
	return nil
}

func (v StreamRequestValidator) Value() entity.StreamRequest {
	return v.output
}

func decryptMagnet(encryptedMagnet string) (string, error) {
	t, err := cipher.NewCryptograph(os.Getenv("AES_KEY"))
	if err != nil {
		return "", err
	}

	plainMagnet, err := hex.DecodeString(encryptedMagnet)
	if err != nil {
		return "", err
	}

	decryptedMagnet, err := t.Decrypt(plainMagnet)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(decryptedMagnet, "magnet:?") {
		return "", ErrInvalidMagnet
	}
	return decryptedMagnet, nil
}
