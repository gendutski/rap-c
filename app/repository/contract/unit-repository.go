package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
)

type UnitRepository interface {
	// create unit
	Create(ctx context.Context, unit *databaseentity.Unit) error
	// get unit by name
	GetUnitByName(ctx context.Context, name string) (*databaseentity.Unit, error)
	// get total units by request param
	GetTotalUnitsByRequest(ctx context.Context, req *payloadentity.GetUnitListRequest) (int64, error)
	// get units by request param
	GetUnitsByRequest(ctx context.Context, req *payloadentity.GetUnitListRequest) ([]*databaseentity.Unit, error)
	// delete unit by name
	Delete(ctx context.Context, unit *databaseentity.Unit) error
}
