package helper

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// generate validator with required struct enabled and json tag
func GenerateStructValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return validate
}
