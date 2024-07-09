package middleware

import (
	"net/http"
	"rap-c/app/entity"

	"github.com/labstack/echo/v4"
)

func PasswordNotChanged(jwtUserContextKey string, isAPI bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get user
			_author := c.Get(jwtUserContextKey)
			user, ok := _author.(*entity.User)
			if !ok || user == nil {
				if isAPI {
					return echo.NewHTTPError(http.StatusUnauthorized)
				} else {
					return c.Redirect(http.StatusFound, entity.WebLoginPath)
				}
			}

			if user.PasswordMustChange {
				if isAPI {
					return &echo.HTTPError{
						Code:     http.StatusForbidden,
						Message:  entity.MiddlewarePasswordNotChangedSamePasswordMessage,
						Internal: entity.NewInternalError(entity.MiddlewarePasswordNotChangedSamePassword, entity.MiddlewarePasswordNotChangedSamePasswordMessage),
					}
				} else {
					return c.Redirect(http.StatusFound, entity.WebPasswordChangePath)
				}
			}

			return next(c)
		}
	}
}
