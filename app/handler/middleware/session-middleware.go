package middleware

import (
	"net/http"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

func ValidateJwtTokenFromSession(sessionUsecase contract.SessionUsecase, route *config.Route, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// create prev route
			sessionUsecase.SetPrevRoute(c)

			//validate token
			user, token, err := sessionUsecase.ValidateJwtToken(c, guestAccepted)
			if err != nil {
				if herr, ok := err.(*echo.HTTPError); ok && herr.Code == http.StatusUnauthorized {
					// set error
					sessionUsecase.SetError(c, herr)
					return c.Redirect(http.StatusFound, route.LoginWebPage.Path())
				}
				return err
			}

			// set context
			c.Set(config.EchoJwtUserContextKey, user)
			c.Set(config.EchoTokenContextKey, token)
			return next(c)
		}
	}
}
