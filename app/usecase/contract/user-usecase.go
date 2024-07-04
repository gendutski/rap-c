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
	// convert user to jwt token
	GenerateJwtToken(ctx context.Context, user *entity.User) (string, error)
	// validate jwt token into user
	ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*entity.User, error)
	// validate jwt token from session
	ValidateSessionJwtToken(ctx context.Context, r *http.Request, w http.ResponseWriter, store sessions.Store, guestAccepted bool) (*entity.User, error)
	// update or modify user password with new password
	RenewPassword(ctx context.Context, user *entity.User, payload *entity.RenewPasswordPayload) error
}
