package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase interface {
	// attempt to login with email and password
	AttemptLogin(ctx context.Context, payload *payloadentity.AttemptLoginPayload) (*databaseentity.User, error)
	// attempt to login with guest account if exists
	AttemptGuestLogin(ctx context.Context) (*databaseentity.User, error)
	// convert user to jwt token
	GenerateJwtToken(ctx context.Context, user *databaseentity.User, isLongSession bool) (string, error)
	// validate jwt token into user
	ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*databaseentity.User, error)
	// update or modify user password with new password
	RenewPassword(ctx context.Context, user *databaseentity.User, payload *payloadentity.RenewPasswordPayload) error
	// submit reset password request
	RequestResetPassword(ctx context.Context, payload *payloadentity.RequestResetPayload) (*databaseentity.User, *databaseentity.PasswordResetToken, error)
	// validate reset password from email
	ValidateResetToken(ctx context.Context, payload *payloadentity.ValidateResetTokenPayload) error
	// submit reset password
	SubmitResetPassword(ctx context.Context, payload *payloadentity.ResetPasswordPayload) (*databaseentity.User, error)
}
