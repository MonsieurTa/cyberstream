package model

import (
	"fmt"
	"reflect"
)

type Params struct {
	field string
	tag   string
}

func GenerateModelGORM(params []Params, model interface{}) interface{} {
	for _, v := range params {
		field, ok := reflect.TypeOf(model).FieldByName(v.field)
		if !ok {
			panic(fmt.Sprintf("field `%s` not found\n", v.field))
		}
		field.Tag = reflect.StructTag(v.tag)
	}
	return model
}
