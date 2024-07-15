package entity

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// get field json value
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	// validate username
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
		return re.MatchString(fl.Field().String())
	})

	return &Validator{
		validate: validate,
	}
}

type Validator struct {
	validate *validator.Validate
}

type ValidatorMessage struct {
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func (e *Validator) Validate(payload interface{}) error {
	err := e.validate.Struct(payload)
	var messages map[string][]*ValidatorMessage

	if err != nil {
		messages = map[string][]*ValidatorMessage{}
		for _, v := range err.(validator.ValidationErrors) {
			messages[v.Field()] = append(messages[v.Field()], &ValidatorMessage{
				Tag:   v.Tag(),
				Param: v.Param(),
			})
		}

		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  messages,
			Internal: NewInternalError(ValidatorBadRequest, ValidatorBadRequestMessage),
		}
	}
	return nil
}
