package responseentity

import (
	payloadentity "rap-c/app/entity/payload-entity"
	"time"
)

type UserResponse struct {
	Username           string    `json:"username"`
	FullName           string    `json:"fullName"`
	Email              string    `json:"email"`
	PasswordMustChange bool      `json:"passwordMustChange"`
	Disabled           bool      `json:"disabled"`
	IsGuest            bool      `json:"isGuest"`
	CreatedAt          time.Time `json:"createdAt"`
	CreatedBy          string    `json:"createdBy"`
	UpdatedAt          time.Time `json:"updatedAt"`
	UpdatedBy          string    `json:"updatedBy"`
}

type GetUserListResponse struct {
	Users   []*UserResponse
	Request *payloadentity.GetUserListRequest
}
