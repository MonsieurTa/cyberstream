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
	"github.com/google/uuid"
)

var (
	ErrInvalidUUID   = errors.New("invalid uuid")
	ErrInvalidMagnet = errors.New("invalid magnet")
)

type UnstoredVideoValidator struct {
	Video struct {
		ID     string `json:"id" binding:"required,uuid"`
		Name   string `json:"name" binding:"required"`
		Magnet string `json:"magnet" binding:"required"`
	} `json:"video"`

	output entity.Video `json:"-"`
}

func NewUnstoredVideoValidator() *UnstoredVideoValidator {
	return &UnstoredVideoValidator{}
}

func (v *UnstoredVideoValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(v.Video.ID)
	if err != nil {
		return ErrInvalidUUID
	}

	decryptedMagnet, err := decryptMagnet(v.Video.Magnet)
	if err != nil {
		return err
	}

	v.output.ID = id
	v.output.Name = v.Video.Name
	v.output.Magnet = decryptedMagnet
	return nil
}

func (v UnstoredVideoValidator) Value() entity.Video {
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
