package contract

import (
	"context"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
)

type AuthRepository interface {
	// attempt login
	DoUserLogin(ctx context.Context, payload *payloadentity.AttemptLoginPayload) (*databaseentity.User, error)
	// renew password
	DoRenewPassword(ctx context.Context, user *databaseentity.User, payload *payloadentity.RenewPasswordPayload) error
	// get user by email
	GetUserByEmail(ctx context.Context, email string) (*databaseentity.User, error)
	// set user reset password
	GenerateUserResetPassword(ctx context.Context, payload *payloadentity.RequestResetPayload) (*databaseentity.PasswordResetToken, error)
	// validate reset password token
	ValidateResetToken(ctx context.Context, payload *payloadentity.ValidateResetTokenPayload) (*databaseentity.PasswordResetToken, error)
	// reset password
	DoResetPassword(ctx context.Context, user *databaseentity.User, reset *databaseentity.PasswordResetToken) error
}
