package common

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CommonError struct {
	Errors map[string]interface{}
}

func MakeCommonError() CommonError {
	return CommonError{
		Errors: make(map[string]interface{}),
	}
}

func NewValidationError(err error) CommonError {
	cerr := MakeCommonError()
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		if v.Param() != "" {
			cerr.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			cerr.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}
	}
	return cerr
}

func NewError(key string, err error) CommonError {
	cerr := MakeCommonError()
	cerr.Errors[key] = err.Error()
	return cerr
}
