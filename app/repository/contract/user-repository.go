package contract

import (
	"context"
	"rap-c/app/entity"
)

type UserRepository interface {
	// create user
	Create(ctx context.Context, user *entity.User) error
	// update existing user
	Update(ctx context.Context, user *entity.User) error
	// get exact user by field: id, username, email
	GetUserByField(ctx context.Context, fieldName string, fieldValue interface{}, notFoundStatus int) (*entity.User, error)
	// get total users by request param
	GetTotalUsersByRequest(ctx context.Context, req *entity.GetUserListRequest) (int64, error)
	// get users by request param
	GetUsersByRequest(ctx context.Context, req *entity.GetUserListRequest) ([]*entity.User, error)
	// set user reset password
	GenerateUserResetPassword(ctx context.Context, email string) (*entity.PasswordResetToken, error)
	// validate reset password token
	ValidateResetToken(ctx context.Context, email string, token string) (*entity.PasswordResetToken, error)
	// reset password
	ResetPassword(ctx context.Context, user *entity.User, reset *entity.PasswordResetToken) error
}
