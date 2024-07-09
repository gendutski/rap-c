package userusecase

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"

	"github.com/labstack/echo/v4"
)

// encrypt password given or generated one
func (uc *usecase) generateUserPassword(pass string) (password string, encryptPassword string, err error) {
	if pass == "" {
		// generate password
		password, err = helper.GenerateStrongPassword()
		if err != nil {
			err = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.UserUsecaseGenerateStrongPasswordError, err.Error()),
			}
			return
		}
	} else {
		password = pass
	}

	// encrypt password
	encryptPassword, err = helper.EncryptPassword(password)
	if err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserUsecaseEncryptPasswordError, err.Error()),
		}
	}
	return
}
