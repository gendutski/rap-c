package route

import (
	"encoding/json"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler/api"
	"rap-c/app/handler/middleware"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	Config      *config.Config
	Route       *config.Route
	AuthUsecase contract.AuthUsecase
	AuthAPI     api.AuthAPI
	UserAPI     api.UserAPI
}

func SetAPIRoute(e *echo.Echo, h *APIHandler) {
	// all login user group
	allLoginRole := []echo.MiddlewareFunc{
		middleware.GetJWT([]byte(h.Config.JwtSecret())),
		middleware.GetUserFromJWT(h.AuthUsecase, h.Config.EnableGuestLogin()),
		middleware.PasswordNotChanged(true, h.Route),
	}

	// non guest only group
	nonGuestOnly := []echo.MiddlewareFunc{
		middleware.GetJWT([]byte(h.Config.JwtSecret())),
		middleware.GetUserFromJWT(h.AuthUsecase, false),
		middleware.PasswordNotChanged(true, h.Route),
	}

	// set api
	h.setAuthAPI(e, nonGuestOnly)
	h.setUserAPI(e, allLoginRole, nonGuestOnly)
}

func (h *APIHandler) setAuthAPI(e *echo.Echo, nonGuestOnly []echo.MiddlewareFunc) {
	// login routes
	e.Add(h.Route.LoginAPI.Method(), h.Route.LoginAPI.Path(), h.AuthAPI.Login)
	// guest login if enabled
	e.Add(h.Route.GuestLoginAPI.Method(), h.Route.GuestLoginAPI.Path(), h.AuthAPI.GuestLogin)
	// renew password routes
	e.Add(h.Route.PasswordMustChangeAPI.Method(), h.Route.PasswordMustChangeAPI.Path(), h.AuthAPI.RenewPassword, nonGuestOnly[:2]...)
	// forgot password (request reset password)
	e.Add(h.Route.RequestResetPasswordAPI.Method(), h.Route.RequestResetPasswordAPI.Path(), h.AuthAPI.RequestResetPassword)
	// reset password
	e.Add(h.Route.ResetPasswordAPI.Method(), h.Route.ResetPasswordAPI.Path(), h.AuthAPI.ResetPassword)
}

func (h *APIHandler) setUserAPI(e *echo.Echo, allLoginRole []echo.MiddlewareFunc, nonGuestOnly []echo.MiddlewareFunc) {
	// all user
	// detail user
	e.Add(h.Route.DetailUserAPI.Method(), h.Route.DetailUserAPI.Path(), h.UserAPI.GetUserDetailByUsername, allLoginRole...)
	// user list
	e.Add(h.Route.ListUserAPI.Method(), h.Route.ListUserAPI.Path(), h.UserAPI.GetUserList, allLoginRole...)
	// user list total
	e.Add(h.Route.TotalUserAPI.Method(), h.Route.TotalUserAPI.Path(), h.UserAPI.GetTotalUserList, allLoginRole...)

	// non guest
	// create new user
	e.Add(h.Route.CreateUserAPI.Method(), h.Route.CreateUserAPI.Path(), h.UserAPI.Create, nonGuestOnly...)
	// update current user
	e.Add(h.Route.UpdateUserAPI.Method(), h.Route.UpdateUserAPI.Path(), h.UserAPI.Update, nonGuestOnly...)
	// update user active status
	e.Add(h.Route.SetStatusUserAPI.Method(), h.Route.SetStatusUserAPI.Path(), echo.NotFoundHandler, nonGuestOnly...)
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
	case map[string][]*entity.ValidatorMessage:
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
