package payloadentity

import "rap-c/app/entity"

// bind struct for get unit list request
type GetUnitListRequest struct {
	Name            string            `query:"name" json:"name"`
	SortField       string            `query:"sortField" json:"sortField"`
	DescendingOrder bool              `query:"descendingOrder" json:"descendingOrder"`
	Limit           int               `query:"limit" json:"limit"`
	Page            entity.Pagination `query:"page" json:"page"`
}

// create unit payload
type CreateUnitPayload struct {
	Name string `json:"name" form:"name" validate:"required,max=30"`
}

// delete unit payload
type DeleteUnitPayload struct {
	Name string `json:"name" form:"name" validate:"required"`
}
