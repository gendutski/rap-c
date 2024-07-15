package route

import (
	"net/http"
	"path/filepath"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"rap-c/app/handler/middleware"
	"rap-c/app/handler/web"

	"github.com/labstack/echo/v4"
)

const (
	storagePath string = "storage"
	assetPath   string = "public-asset"
	imagePath   string = "images"
	favIcon     string = "favicon.ico"
)

type WebHandler struct {
	Config         *config.Config
	Route          *config.Route
	AuthUsecase    contract.AuthUsecase
	SessionUsecase contract.SessionUsecase
	AuthPage       web.AuthPage
	UserPage       web.UserPage
	DashboardPage  web.DashboardPage
}

func SetWebRoute(e *echo.Echo, h *WebHandler) {
	// home page
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, h.Route.LoginWebPage.Path())
	})
	// asset folder
	e.Static("/assets", filepath.Join(storagePath, assetPath))
	e.File("/favicon.ico", filepath.Join(storagePath, assetPath, imagePath, favIcon))

	// all login user group
	allLoginRole := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.SessionUsecase, h.Route, h.Config.EnableGuestLogin()),
		middleware.PasswordNotChanged(false, h.Route),
	}

	// non guest only group
	nonGuestOnly := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.SessionUsecase, h.Route, false),
		middleware.PasswordNotChanged(false, h.Route),
	}

	// user api
	h.setAuthWebPage(e)
	h.setMainWebPage(e, allLoginRole, nonGuestOnly)
}

func (h *WebHandler) setAuthWebPage(e *echo.Echo) {
	// login page
	e.Add(h.Route.LoginWebPage.Method(), h.Route.LoginWebPage.Path(), h.AuthPage.Login)
	// token session submit page
	e.Add(h.Route.SubmitTokenSessionWebPage.Method(), h.Route.SubmitTokenSessionWebPage.Path(), h.AuthPage.SubmitToken)
	// logout
	e.Add(h.Route.LogoutWebPage.Method(), h.Route.LogoutWebPage.Path(), h.AuthPage.Logout)
	// forgot password
	e.Add(h.Route.ForgotPasswordWebPage.Method(), h.Route.ForgotPasswordWebPage.Path(), h.AuthPage.ForgotPassword)

	// // reset password
	// e.GET(entity.WebRequestResetPath, h.AuthPage.RequestResetPassword)
	// e.GET(entity.WebResetPasswordPath, h.AuthPage.ResetPassword)

	// password must change
	e.Add(h.Route.PasswordMustChangeWebPage.Method(), h.Route.PasswordMustChangeWebPage.Path(),
		h.AuthPage.PasswordMustChange, middleware.ValidateJwtTokenFromSession(h.SessionUsecase, h.Route, false))
}

func (h *WebHandler) setMainWebPage(e *echo.Echo, allLoginRole []echo.MiddlewareFunc, nonGuestOnly []echo.MiddlewareFunc) {
	// all user
	// profile page
	e.Add(h.Route.ProfileWebPage.Method(), h.Route.ProfileWebPage.Path(), h.UserPage.Profile, allLoginRole...)
	// dashboard
	e.Add(h.Route.DashboardWebPage.Method(), h.Route.DashboardWebPage.Path(), h.DashboardPage.Dashboard, allLoginRole...)

	// // non guest
	// e.GET("/user", echo.NotFoundHandler, nonGuestOnly...)
}

func WebErrorHandler(e *echo.Echo, err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if report.Code == http.StatusNotFound {
		c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
			"code":    report.Code,
			"title":   http.StatusText(report.Code),
			"message": report.Message,
		})
		return
	} else if report.Code == http.StatusUnauthorized {
		c.Render(http.StatusUnauthorized, "401.html", nil)
		return
	}
	c.Render(report.Code, "error.html", map[string]interface{}{
		"code":    report.Code,
		"title":   http.StatusText(report.Code),
		"message": report.Message,
	})
}
