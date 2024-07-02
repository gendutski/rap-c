package middleware

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	usermodule "rap-c/app/module/user-module"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func ValidateJwtTokenFromSession(store sessions.Store, jwtSecret []byte, jwtUserContextKey string, userModule usermodule.UserUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get token from session
			sess, err := helper.NewSession(c.Request(), c.Response(), store, entity.SessionID)
			if err != nil {
				return err
			}
			tokenStr, ok := sess.Get(entity.TokenSessionName).(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			// parse token
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			// validate user token
			ctx := c.Request().Context()
			user, err := userModule.ValidateJwtToken(ctx, token, guestAccepted)
			if err != nil {
				return err
			}
			c.Set(jwtUserContextKey, user)
			return next(c)
		}
	}
}
