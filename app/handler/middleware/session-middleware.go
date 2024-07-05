package middleware

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/usecase/contract"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func ValidateJwtTokenFromSession(store sessions.Store, jwtUserContextKey string, userUsecase contract.UserUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			user, token, err := userUsecase.ValidateSessionJwtToken(ctx, c.Request(), c.Response(), store, guestAccepted)
			if err != nil {
				if herr, ok := err.(*echo.HTTPError); ok && herr.Code == http.StatusUnauthorized {
					return c.Redirect(http.StatusMovedPermanently, entity.WebLoginPath)
				}
				return err
			}
			c.Set(jwtUserContextKey, user)
			c.Set(entity.TokenSessionName, token)
			return next(c)
		}
	}
}
