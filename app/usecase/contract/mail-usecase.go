package contract

import (
	databaseentity "rap-c/app/entity/database-entity"
)

type MailUsecase interface {
	Welcome(user *databaseentity.User, password string) error
	ResetPassword(user *databaseentity.User, token *databaseentity.PasswordResetToken) error
	UpdateUser(user *databaseentity.User) error
	UpdateActiveStatusUser(user *databaseentity.User) error
}
