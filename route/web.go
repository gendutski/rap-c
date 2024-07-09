package route

import (
	"net/http"
	"path/filepath"
	"rap-c/app/entity"
	"rap-c/app/usecase/contract"

	"rap-c/app/handler/middleware"
	"rap-c/app/handler/web"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	storagePath string = "storage"
	assetPath   string = "public-asset"
	imagePath   string = "images"
	favIcon     string = "favicon.ico"
)

type WebHandler struct {
	JwtUserContextKey string
	JwtSecret         string
	GuestAccepted     bool
	AuthUsecase       contract.AuthUsecase
	Store             sessions.Store
	AuthPage          web.AuthPage
	UserPage          web.UserPage
}

func SetWebRoute(e *echo.Echo, h *WebHandler) {
	// home page
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/login")
	})
	// asset folder
	e.Static("/assets", filepath.Join(storagePath, assetPath))
	e.File("/favicon.ico", filepath.Join(storagePath, assetPath, imagePath, favIcon))

	// all login user group
	allLoginRole := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.AuthUsecase, h.GuestAccepted),
		middleware.PasswordNotChanged(h.JwtUserContextKey, false),
	}

	// non guest only group
	nonGuestOnly := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.AuthUsecase, false),
		middleware.PasswordNotChanged(h.JwtUserContextKey, false),
	}

	// user api
	h.setAuthWebPage(e)
	h.setUserWebPage(e, allLoginRole, nonGuestOnly)
}

func (h *WebHandler) setAuthWebPage(e *echo.Echo) {
	// login page
	e.GET(entity.WebLoginPath, h.AuthPage.Login).Name = entity.LoginRouteName
	e.POST("/guest-login", h.AuthPage.GuestLogin).Name = entity.GuestLoginRouteName
	e.POST("/submit-login", h.AuthPage.PostLogin).Name = entity.PostLoginRouteName
	e.POST(entity.WebLogoutPath, h.AuthPage.PostLogout).Name = entity.PostLogoutRouteName

	// reset password
	e.GET("/request-reset", h.AuthPage.RequestResetPassword).Name = entity.RequestResetPasswordName
	e.POST("/request-reset", h.AuthPage.SubmitRequestResetPassword).Name = entity.PostRequestResetPasswordName
	e.GET(entity.WebResetPasswordPath, h.AuthPage.ResetPassword).Name = entity.ResetPasswordName
	e.POST(entity.WebResetPasswordPath, h.AuthPage.SubmitResetPassword).Name = entity.SubmitResetPasswordName

	// password must change
	e.GET(entity.WebPasswordChangePath, h.AuthPage.PasswordChanger, middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.AuthUsecase, h.GuestAccepted))
	e.POST(entity.WebPasswordChangePath, h.AuthPage.SubmitPasswordChanger, middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.AuthUsecase, h.GuestAccepted)).Name = entity.RenewPasswordRouteName

}

func (h *WebHandler) setUserWebPage(e *echo.Echo, allLoginRole []echo.MiddlewareFunc, nonGuestOnly []echo.MiddlewareFunc) {
	// all user
	e.GET("/profile", h.UserPage.Profile, allLoginRole...).Name = entity.ProfileRouteName

	// non guest
	e.GET("/user", echo.NotFoundHandler, nonGuestOnly...)
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
