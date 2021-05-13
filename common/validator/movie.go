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

type UnstoredMovieValidator struct {
	Movie struct {
		ID     string `json:"id" binding:"required,uuid"`
		Name   string `json:"name" binding:"required"`
		Magnet string `json:"magnet" binding:"required"`
	} `json:"movie"`

	output entity.Movie `json:"-"`
}

func NewUnstoredMovieValidator() *UnstoredMovieValidator {
	return &UnstoredMovieValidator{}
}

func (v *UnstoredMovieValidator) Validate(c *gin.Context) error {
	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(v.Movie.ID)
	if err != nil {
		return ErrInvalidUUID
	}

	decryptedMagnet, err := decryptMagnet(v.Movie.Magnet)
	if err != nil {
		return err
	}

	v.output.ID = id
	v.output.Name = v.Movie.Name
	v.output.Magnet = decryptedMagnet
	return nil
}

func (v UnstoredMovieValidator) Value() entity.Movie {
	return v.output
}

func decryptMagnet(encryptedMagnet string) (string, error) {
	t, err := cipher.NewTranslator(os.Getenv("AES_KEY"))
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
