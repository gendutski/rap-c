package middleware

import (
	"errors"
	usermodule "rap-c/app/module/user-module"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	jwtContextKey string = "userToken"
)

func GetJWT(jwtSecret []byte) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: jwtSecret,
		ContextKey: jwtContextKey,
	})
}

func GetUserFromJWT(jwtUserContextKey string, userModule usermodule.UserUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// validate token with user db middleware
		return func(c echo.Context) error {
			token, ok := c.Get(jwtContextKey).(*jwt.Token) // by default token is stored under `user` key
			if !ok {
				c.Error(errors.New("JWT token missing or invalid"))
			}

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
