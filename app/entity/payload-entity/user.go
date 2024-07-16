package payloadentity

import (
	"rap-c/app/entity"
)

type GetUserDetailRequest struct {
	Username string `param:"username" json:"username" validate:"required"`
}

// bind struct for get user list request
type GetUserListRequest struct {
	UserName        string            `query:"username" json:"username"`
	Email           string            `query:"email" json:"email"`
	FullName        string            `query:"fullName" json:"fullName"`
	Show            string            `query:"show" json:"show"`
	SortField       string            `query:"sortField" json:"sortField"`
	DescendingOrder bool              `query:"descendingOrder" json:"descendingOrder"`
	GuestOnly       bool              `query:"-" json:"-"`
	Limit           int               `query:"limit" json:"limit"`
	Page            entity.Pagination `query:"page" json:"page"`
}

// create user payload
type CreateUserPayload struct {
	Username string `json:"username" form:"username" validate:"required,max=30,username"`
	FullName string `json:"fullName" form:"fullName" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	IsGuest  bool   `json:"-"`
}

// update user payload
type UpdateUserPayload struct {
	Username        string `json:"username" form:"username" validate:"omitempty,max=30,username"`
	FullName        string `json:"fullName" form:"fullName"`
	Email           string `json:"email" form:"email" validate:"omitempty,email"`
	Password        string `json:"password" form:"password" validate:"omitempty,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required_with=Password,eqfield=Password"`
}
