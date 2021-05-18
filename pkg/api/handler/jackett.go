package handler

import (
	"net/http"

	"github.com/MonsieurTa/hypertube/common/validator"
	"github.com/MonsieurTa/hypertube/pkg/api/usecase/jackett"
	"github.com/gin-gonic/gin"
)

func JackettSearch(service jackett.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := validator.NewSearchValidator()
		err := validator.Validate(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pattern, categories := validator.Value()

		res, err := service.Search(pattern, categories)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func JackettCategories(service jackett.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats := service.Categories()
		c.JSON(http.StatusOK, cats)
	}
}
