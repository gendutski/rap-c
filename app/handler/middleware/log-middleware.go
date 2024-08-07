package middleware

import (
	"rap-c/app/entity"
	"rap-c/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetLog(logMode config.LogMode, enableWarnFileLog bool) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			entity.InitLog(v.URI, c.Request().Method, "request", v.Status, v.Error, logMode, enableWarnFileLog).Log()
			return nil
		},
	})
}
