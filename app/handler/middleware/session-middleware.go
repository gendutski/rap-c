package middleware

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

func ValidateJwtTokenFromSession(sessionUsecase contract.SessionUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, token, err := sessionUsecase.ValidateJwtToken(c, guestAccepted)
			if err != nil {
				if herr, ok := err.(*echo.HTTPError); ok && herr.Code == http.StatusUnauthorized {
					sessionUsecase.SetError(c, herr)
					return c.Redirect(http.StatusFound, entity.WebLoginPath)
				}
				return err
			}
			c.Set(config.EchoJwtUserContextKey, user)
			c.Set(config.EchoTokenContextKey, token)
			return next(c)
		}
	}
}
