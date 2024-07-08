package contract

import (
	"context"
	"net/http"
	"rap-c/app/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

type UserUsecase interface {
	// create user
	Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, string, error)
	// attempt to login with email and password
	AttemptLogin(ctx context.Context, payload *entity.AttemptLoginPayload) (*entity.User, error)
	// attempt to login with guest account if exists
	AttemptGuestLogin(ctx context.Context) (*entity.User, error)
	// convert user to jwt token
	GenerateJwtToken(ctx context.Context, user *entity.User, isLongSession bool) (string, error)
	// validate jwt token into user
	ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*entity.User, error)
	// validate jwt token from session
	ValidateSessionJwtToken(ctx context.Context, r *http.Request, w http.ResponseWriter, store sessions.Store, guestAccepted bool) (*entity.User, string, error)
	// update or modify user password with new password
	RenewPassword(ctx context.Context, user *entity.User, payload *entity.RenewPasswordPayload) error
	// get user list
	GetUserList(ctx context.Context, req *entity.GetUserListRequest) ([]*entity.User, error)
	// get total user list
	GetTotalUserList(ctx context.Context, req *entity.GetUserListRequest) (int64, error)
	// get user by username
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	// submit reset password request
	RequestResetPassword(ctx context.Context, email string) (*entity.User, *entity.PasswordResetToken, error)
	// subit reset password
	ValidateResetPassword(ctx context.Context, email string, token string) error
}
