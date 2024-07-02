package middleware

import (
	usermodule "rap-c/app/module/user-module"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func ValidateJwtTokenFromSession(store sessions.Store, jwtUserContextKey string, userModule usermodule.UserUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			user, err := userModule.ValidateSessionJwtToken(ctx, c.Request(), c.Response(), store, guestAccepted)
			if err != nil {
				return err
			}
			c.Set(jwtUserContextKey, user)
			return next(c)
		}
	}
}
