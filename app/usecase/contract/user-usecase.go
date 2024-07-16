package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
)

type UserUsecase interface {
	// create user
	Create(ctx context.Context, payload *payloadentity.CreateUserPayload, author *databaseentity.User) (*databaseentity.User, string, error)
	// get user list
	GetUserList(ctx context.Context, req *payloadentity.GetUserListRequest) ([]*databaseentity.User, error)
	// get total user list
	GetTotalUserList(ctx context.Context, req *payloadentity.GetUserListRequest) (int64, error)
	// get user by username
	GetUserByUsername(ctx context.Context, req *payloadentity.GetUserDetailRequest) (*databaseentity.User, error)
	// update current user data
	Update(ctx context.Context, payload *payloadentity.UpdateUserPayload, author *databaseentity.User) error
}
