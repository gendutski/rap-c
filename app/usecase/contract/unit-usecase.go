package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
)

type UnitUsecase interface {
	// create unit
	Create(ctx context.Context, payload *payloadentity.CreateUnitPayload, author *databaseentity.User) (*databaseentity.Unit, error)
	// get unit list
	GetUnitList(ctx context.Context, req *payloadentity.GetUnitListRequest) ([]*databaseentity.Unit, error)
	// get total unit list
	GetTotalUnitList(ctx context.Context, req *payloadentity.GetUnitListRequest) (int64, error)
	// delete unit
	Delete(ctx context.Context, payload *payloadentity.DeleteUnitPayload) error
}
