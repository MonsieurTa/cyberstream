package validator

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(e *gin.Engine) bool {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phone_fr", phoneValidation)
		return true
	}
	return false
}

var phoneValidation validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	regexp := regexp.MustCompile(`^(?:(?:\+|00)33|0)\s*[1-9](?:[\s.-]*\d{2}){4}$`)
	return regexp.MatchString(phone)
}
