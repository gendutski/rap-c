package responseentity

import (
	payloadentity "rap-c/app/entity/payload-entity"
	"time"
)

type UnitResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type GetUnitListResponse struct {
	Units   []*UnitResponse
	Request *payloadentity.GetUnitListRequest
}
