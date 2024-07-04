package route

import (
	"encoding/json"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler/api"
	"rap-c/app/handler/middleware"
	"rap-c/app/usecase/contract"

	"github.com/labstack/echo/v4"
)

const (
	ApiGroup string = "/api"
)

type APIHandler struct {
	JwtUserContextKey string
	JwtSecret         string
	GuestAccepted     bool
	UserUsecase       contract.UserUsecase
	UserAPI           api.UserAPI
}

func SetAPIRoute(e *echo.Echo, h *APIHandler) {
	// api group
	apiGroup := e.Group(ApiGroup)

	// non login routes
	apiGroup.POST("/login", h.UserAPI.Login)

	// renew password routes
	apiGroup.PUT("/renew-password", h.UserAPI.RenewPassword,
		middleware.GetJWT([]byte(h.JwtSecret)),
		middleware.GetUserFromJWT(h.JwtUserContextKey, h.UserUsecase, false))

	// all login user group
	allLoginRole := []echo.MiddlewareFunc{
		middleware.GetJWT([]byte(h.JwtSecret)),
		middleware.GetUserFromJWT(h.JwtUserContextKey, h.UserUsecase, h.GuestAccepted),
		middleware.PasswordNotChanged(h.JwtUserContextKey, true),
	}

	// non guest only group
	nonGuestOnly := []echo.MiddlewareFunc{
		middleware.GetJWT([]byte(h.JwtSecret)),
		middleware.GetUserFromJWT(h.JwtUserContextKey, h.UserUsecase, false),
		middleware.PasswordNotChanged(h.JwtUserContextKey, true),
	}

	// user api
	h.setUserAPI(apiGroup, allLoginRole, nonGuestOnly)
}

func (h *APIHandler) setUserAPI(apiGroup *echo.Group, allLoginRole []echo.MiddlewareFunc, nonGuestOnly []echo.MiddlewareFunc) {
	// all user
	apiGroup.GET("/user/detail/:username", echo.NotFoundHandler, allLoginRole...)
	apiGroup.GET("/user/list", echo.NotFoundHandler, allLoginRole...)
	apiGroup.GET("/user/total", echo.NotFoundHandler, allLoginRole...)

	// non guest
	apiGroup.POST("/user/create", h.UserAPI.Create, nonGuestOnly...)
	apiGroup.PUT("/user/update", echo.NotFoundHandler, nonGuestOnly...)
	apiGroup.PUT("/user/active-status", echo.NotFoundHandler, nonGuestOnly...)
}

// error handler
func APIErrorHandler(e *echo.Echo, err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var code int
	errMessage := err.Error()
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			} else if herr, ok := he.Internal.(*entity.InternalError); ok {
				code = herr.Code
				errMessage = herr.Message
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	// Issue #1426
	message := he.Message

	switch m := he.Message.(type) {
	case string:
		if e.Debug {
			message = echo.Map{"code": code, "message": m, "error": errMessage}
		} else {
			message = echo.Map{"code": code, "message": m}
		}
	case []string:
		if e.Debug {
			message = echo.Map{"code": code, "message": m, "error": errMessage}
		} else {
			message = echo.Map{"code": code, "message": m}
		}
	case json.Marshaler:
		// do nothing - this type knows how to format itself to JSON
	case error:
		message = echo.Map{"code": code, "message": m.Error()}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(he.Code, message)
	}
	if err != nil {
		e.Logger.Error(err)
	}
}
