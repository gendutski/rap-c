package contract

import "rap-c/app/entity"

type MailUsecase interface {
	Welcome(user *entity.User, password string) error
}