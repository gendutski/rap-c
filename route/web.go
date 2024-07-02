package route

import (
	"net/http"
	"path/filepath"
	"rap-c/app/entity"

	"rap-c/app/handler/middleware"
	"rap-c/app/handler/web"
	usermodule "rap-c/app/module/user-module"

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
	UserModule        usermodule.UserUsecase
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
	e.GET("/login", h.UserPage.Login).Name = entity.LoginRouteName
	e.POST("/submit-login", h.UserPage.PostLogin).Name = entity.PostLoginRouteName
	e.POST("/submit-logout", h.UserPage.PostLogout).Name = entity.PostLogoutRouteName

	// all login user group
	allLoginRole := e.Group(
		"",
		middleware.ValidateJwtTokenFromSession(h.Store, []byte(h.JwtSecret), h.JwtUserContextKey, h.UserModule, h.GuestAccepted),
	)

	// non guest only group
	nonGuestOnly := e.Group(
		"",
		middleware.ValidateJwtTokenFromSession(h.Store, []byte(h.JwtSecret), h.JwtUserContextKey, h.UserModule, false),
	)

	// user api
	h.setUserWebPage(allLoginRole, nonGuestOnly)
}

func (h *WebHandler) setUserWebPage(allLoginRole *echo.Group, nonGuestOnly *echo.Group) {
	// all user
	allLoginRole.GET("/profile", h.UserPage.Profile).Name = entity.ProfileRouteName

	// non guest
	nonGuestOnly.GET("/user", echo.NotFoundHandler)
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
	}
	c.Render(http.StatusInternalServerError, "500.html", map[string]interface{}{"message": report.Message})
}
