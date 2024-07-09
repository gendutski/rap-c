package contract

import (
	"context"
	"rap-c/app/entity"
)

type UserUsecase interface {
	// create user
	Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, string, error)
	// get user list
	GetUserList(ctx context.Context, req *entity.GetUserListRequest) ([]*entity.User, error)
	// get total user list
	GetTotalUserList(ctx context.Context, req *entity.GetUserListRequest) (int64, error)
	// get user by username
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}
