package entity

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// attempt to login payload
type AttemptLoginPayload struct {
	Email      string `json:"email" form:"email" validate:"required,email"`
	Password   string `json:"password" form:"password" validate:"required"`
	RememberMe bool   `json:"rememberMe" form:"rememberMe"`
}

// will validate payload, and return slice of error messages or nil
func (e AttemptLoginPayload) Validate(validate *validator.Validate) []string {
	err := validate.Struct(e)
	var messages []string
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			switch v.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("field `%s` is required", v.Field()))
			case "email":
				messages = append(messages, "invalid email")
			}
		}
	}
	return messages
}

type RenewPasswordPayload struct {
	Password        string `json:"password" form:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

// will validate payload, and return slice of error messages or nil
func (e RenewPasswordPayload) Validate(validate *validator.Validate) []string {
	err := validate.Struct(e)
	var messages []string
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			switch v.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("field `%s` is required", v.Field()))
			case "eqfield":
				messages = append(messages, "password confirmation is not same")
			}
		}
	}
	return messages
}

type ResetPasswordPayload struct {
	Email           string `json:"email" form:"email" validate:"required,email"`
	Token           string `json:"token" form:"token" validate:"required"`
	Password        string `json:"password" form:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

// will validate payload, and return slice of error messages or nil
func (e ResetPasswordPayload) Validate(validate *validator.Validate) []string {
	err := validate.Struct(e)
	var messages []string
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			switch v.Tag() {
			case "email":
				messages = append(messages, "invalid email")
			case "required":
				messages = append(messages, fmt.Sprintf("field `%s` is required", v.Field()))
			case "eqfield":
				messages = append(messages, "password confirmation is not same")
			}
		}
	}
	return messages
}
