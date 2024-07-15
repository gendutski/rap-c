package payloadentity

import "rap-c/app/entity"

// bind struct for get user list request
type GetUserListRequest struct {
	UserName        string            `query:"userName" json:"userName"`
	Email           string            `query:"email" json:"email"`
	FullName        string            `query:"fullName" json:"fullName"`
	Show            string            `query:"show" json:"show"`
	SortField       string            `query:"sortField" json:"sortField"`
	DescendingOrder bool              `query:"descendingOrder" json:"descendingOrder"`
	GuestOnly       bool              `query:"-" json:"-"`
	Limit           int               `query:"limit" json:"limit"`
	Page            entity.Pagination `query:"page" json:"page"`
}
