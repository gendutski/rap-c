package payloadentity

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// create user payload
type CreateUserPayload struct {
	Username string `json:"username" form:"username" validate:"required,max=30,username"`
	FullName string `json:"fullName" form:"fullName" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	IsGuest  bool   `json:"-"`
}

// will validate payload, and return slice of error messages or nil
func (e CreateUserPayload) Validate(validate *validator.Validate) []string {
	err := validate.Struct(e)
	var messages []string
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			switch v.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("field `%s` is required", v.Field()))
			case "email":
				messages = append(messages, "invalid email")
			case "max":
				messages = append(messages, fmt.Sprintf("only allowed max `%s` characters for `%s`", v.Param(), v.Field()))
			case "username":
				messages = append(messages, fmt.Sprintf("only allowed alphanumeric, period, dash, and underscore for `%s`", v.Field()))
			}
		}
	}
	return messages
}
