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
	UserUsecase       contract.UserUsecase
	UserPage          web.UserPage
	Store             sessions.Store
}

func SetWebRoute(e *echo.Echo, h *WebHandler) {
	// home page
	e.GET("/", func(c echo.Context) error {

		return c.JSON(200, c.Echo().Routes())
	})
	// asset folder
	e.Static("/assets", filepath.Join(storagePath, assetPath))
	e.File("/favicon.ico", filepath.Join(storagePath, assetPath, imagePath, favIcon))

	// login page
	e.GET(entity.WebLoginPath, h.UserPage.Login).Name = entity.LoginRouteName
	e.POST("/submit-login", h.UserPage.PostLogin).Name = entity.PostLoginRouteName
	e.POST("/submit-logout", h.UserPage.PostLogout).Name = entity.PostLogoutRouteName

	// password must change
	e.GET(entity.WebPasswordChangePath, h.UserPage.PasswordChanger, middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.UserUsecase, h.GuestAccepted))

	// all login user group
	allLoginRole := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.UserUsecase, h.GuestAccepted),
		middleware.PasswordNotChanged(h.JwtUserContextKey, false),
	}

	// non guest only group
	nonGuestOnly := []echo.MiddlewareFunc{
		middleware.ValidateJwtTokenFromSession(h.Store, h.JwtUserContextKey, h.UserUsecase, false),
		middleware.PasswordNotChanged(h.JwtUserContextKey, false),
	}

	// user api
	h.setUserWebPage(e, allLoginRole, nonGuestOnly)
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
		c.Render(http.StatusNotFound, "404.html", nil)
		return
	} else if report.Code == http.StatusUnauthorized {
		c.Render(http.StatusUnauthorized, "401.html", nil)
		return
	}
	c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
		"code":    report.Code,
		"title":   http.StatusText(report.Code),
		"message": report.Message,
	})
}
