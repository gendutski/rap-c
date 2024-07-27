package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	responseentity "rap-c/app/entity/response-entity"
)

// API response formatter
type FormatterUsecase interface {
	FormatUser(ctx context.Context, user *databaseentity.User, mapUsers map[int]string) (*responseentity.UserResponse, error)
	FormatUsers(ctx context.Context, users []*databaseentity.User) ([]*responseentity.UserResponse, error)
	FormatUnit(ctx context.Context, unit *databaseentity.Unit, mapUsers map[int]string) (*responseentity.UnitResponse, error)
	FormatUnits(ctx context.Context, units []*databaseentity.Unit) ([]*responseentity.UnitResponse, error)
}
