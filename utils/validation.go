package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}

	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			field, ok := reflect.TypeOf(data).FieldByName(v.StructField())
			if !ok {
				continue
			}
			jsonTag := field.Tag.Get("json")
			res[jsonTag] = TranslateTag(v, jsonTag)
		}
	}

	return res
}

func TranslateTag(fd validator.FieldError, jsonTag string) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("Field %s is required", jsonTag)
	case "email":
		return fmt.Sprintf("Field %s must be a valid email address", jsonTag)
	case "min":
		return fmt.Sprintf("Field %s must be at least %s characters long", jsonTag, fd.Param())
	case "max":
		return fmt.Sprintf("Field %s must be at most %s characters long", jsonTag, fd.Param())
	case "len":
		return fmt.Sprintf("Field %s must be exactly %s characters long", jsonTag, fd.Param())
	case "numeric":
		return fmt.Sprintf("Field %s must be a number", jsonTag)
	case "alphanum":
		return fmt.Sprintf("Field %s must contain only letters and numbers", jsonTag)
	case "alpha":
		return fmt.Sprintf("Field %s must contain only letters", jsonTag)
	case "oneof":
		return fmt.Sprintf("Field %s must be one of: %s", jsonTag, fd.Param())
	// More validation tags as needed
	default:
		return "Validation failed " + fd.StructField()
	}
}
