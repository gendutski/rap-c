package contract

import (
	databaseentity "rap-c/app/entity/database-entity"

	"github.com/labstack/echo/v4"
)

type SessionUsecase interface {
	// save jwt token to session
	SaveJwtToken(e echo.Context, token string) error
	// validate jwt token from session
	ValidateJwtToken(e echo.Context, guestAccepted bool) (*databaseentity.User, string, error)
	// set error message to session
	SetError(e echo.Context, err error) error
	// get error message from session
	GetError(e echo.Context) *echo.HTTPError
	// set info message to session
	SetInfo(e echo.Context, info interface{}) error
	// get info message to session
	GetInfo(e echo.Context) (interface{}, error)
	// set previous route (method & path) to session
	SetPrevRoute(e echo.Context) error
	// get previous route (method & path) from session
	GetPrevRoute(e echo.Context) (method string, path string)
	// destroy session
	Logout(e echo.Context) error
}
