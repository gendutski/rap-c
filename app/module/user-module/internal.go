package usermodule

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"

	"github.com/labstack/echo/v4"
)

// encrypt password given or generated one
func (uc *usecase) generateUserPassword(pass string) (string, error) {
	var err error
	if pass == "" {
		// generate password
		pass, err = helper.GenerateStrongPassword()
		if err != nil {
			return "", &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.GeneratePasswordError, err.Error()),
			}
		}
	}

	// encrypt password
	return helper.EncryptPassword(pass)
}
