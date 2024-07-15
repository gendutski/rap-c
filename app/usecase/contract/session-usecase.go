package contract

import (
	databaseentity "rap-c/app/entity/database-entity"

	"github.com/labstack/echo/v4"
)

type SessionUsecase interface {
	SaveJwtToken(e echo.Context, token string) error
	ValidateJwtToken(e echo.Context, guestAccepted bool) (*databaseentity.User, string, error)
	SetError(e echo.Context, err error) error
	GetError(e echo.Context) *echo.HTTPError
	SetInfo(e echo.Context, info interface{}) error
	GetInfo(e echo.Context) (interface{}, error)
	Logout(e echo.Context) error
}
