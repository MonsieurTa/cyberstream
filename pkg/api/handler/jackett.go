package handler

import (
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/common/cipher"
	"github.com/MonsieurTa/hypertube/common/validator"
	jackettService "github.com/MonsieurTa/hypertube/pkg/api/usecase/jackett"
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/go-jackett"
)

func JackettSearch(service jackettService.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewSearchValidator()
		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pattern, categories := validator.Value()

		resp, err := service.Search(pattern, categories)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rv, err := encryptResultsMagnets(resp.Results)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rv)
	}
}

func JackettCategories(service jackettService.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats := service.Categories()
		c.JSON(http.StatusOK, cats)
	}
}

func JackettIndexers(service jackettService.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		idxs, err := service.ConfiguredIndexers()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, idxs)
	}
}

func encryptResultsMagnets(results []jackett.Result) ([]jackett.Result, error) {
	rv := make([]jackett.Result, 0, len(results))

	c, err := cipher.NewCryptograph(os.Getenv("AES_KEY"))
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		v.MagnetUri, err = c.Encrypt([]byte(v.MagnetUri))
		if err != nil {
			return nil, err
		}
		rv = append(rv, v)
	}
	return rv, nil
}
