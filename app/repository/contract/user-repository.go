package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
)

type UserRepository interface {
	// create user
	Create(ctx context.Context, user *databaseentity.User) error
	// update existing user
	Update(ctx context.Context, user *databaseentity.User) error
	// get exact user by field: id, username, email
	GetUserByField(ctx context.Context, fieldName string, fieldValue interface{}, notFoundStatus int) (*databaseentity.User, error)
	// get total users by request param
	GetTotalUsersByRequest(ctx context.Context, req *payloadentity.GetUserListRequest) (int64, error)
	// get users by request param
	GetUsersByRequest(ctx context.Context, req *payloadentity.GetUserListRequest) ([]*databaseentity.User, error)
	// map user username by id
	MapUserUsername(ctx context.Context, userIDs interface{}) (map[int]string, error)
}
