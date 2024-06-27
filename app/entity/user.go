package entity

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	RequestShowActive    string = "active"
	RequestShowNotActive string = "not active"
)

// table users model
type User struct {
	ID                 int       `gorm:"primaryKey" json:"-"`
	Username           string    `gorm:"unique;size:30;not null" json:"userName"`
	FullName           string    `gorm:"size:100;not null" json:"fullName"`
	Email              string    `gorm:"unique;size:100;not null" json:"email"`
	Password           string    `gorm:"size:255;not null" json:"-"`
	PasswordMustChange bool      `gorm:"not null;default:0" json:"passwordMustChange"`
	Disabled           bool      `gorm:"not null;default:0" json:"disabled"`
	IsGuest            bool      `gorm:"not null;default:0" json:"isGuest"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy          string    `gorm:"size:30;not null;default:'SYSTEM'" json:"createdBy"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy          string    `gorm:"size:30;not null;default:'SYSTEM'" json:"updatedBy"`
}

// create user payload
type CreateUserPayload struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
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
			}
		}
	}
	return messages
}

// attempt to login payload
type AttemptLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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

// bind struct for get user list request
type GetUserListRequest struct {
	UserName        string     `query:"userName" json:"userName"`
	Email           string     `query:"email" json:"email"`
	FullName        string     `query:"fullName" json:"fullName"`
	Role            string     `query:"role" json:"role"`
	Show            string     `query:"show" json:"show"`
	SortField       string     `query:"sortField" json:"sortField"`
	DescendingOrder bool       `query:"descendingOrder" json:"descendingOrder"`
	Limit           int        `query:"limit" json:"limit"`
	Page            Pagination `query:"page" json:"page"`
}
