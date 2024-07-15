package middleware

import (
	"errors"
	"rap-c/app/usecase/contract"
	"rap-c/config"

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

func GetUserFromJWT(authUsecase contract.AuthUsecase, guestAccepted bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// validate token with user db middleware
		return func(c echo.Context) error {
			token, ok := c.Get(jwtContextKey).(*jwt.Token) // by default token is stored under `user` key
			if !ok {
				c.Error(errors.New("JWT token missing or invalid"))
			}

			ctx := c.Request().Context()
			user, err := authUsecase.ValidateJwtToken(ctx, token, guestAccepted)
			if err != nil {
				return err
			}
			c.Set(config.EchoJwtUserContextKey, user)
			return next(c)
		}
	}
}
