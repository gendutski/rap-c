package contract

import (
	"context"
	"rap-c/app/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
}
