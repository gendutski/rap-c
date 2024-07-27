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
}
