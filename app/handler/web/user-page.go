package web

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	usermodule "rap-c/app/module/user-module"
	"rap-c/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	loginSessionName string = "logsession"
)

type UserPage interface {
	// login page
	Login(e echo.Context) error
	// post login form
	PostLogin(e echo.Context) error
	// post logout
	PostLogout(e echo.Context) error
	// profile page
	Profile(e echo.Context) error
}

func NewUserPage(cfg config.Config, store sessions.Store, userModule usermodule.UserUsecase) UserPage {
	return &userHandler{cfg, store, userModule}
}

type userHandler struct {
	cfg        config.Config
	store      sessions.Store
	userModule usermodule.UserUsecase
}

func (h *userHandler) Login(e echo.Context) error {
	sess, err := helper.NewSession(e.Request(), e.Response(), h.store, loginSessionName)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  entity.SessionErrorMessage,
			Internal: entity.NewInternalError(entity.SessionError, err.Error()),
		}
	}

	routeMap := helper.RouteMap(e.Echo().Routes())
	var emailValue string

	if _val, ok := sess.Flash("email").(string); ok {
		emailValue = _val
	}

	return e.Render(http.StatusOK, "login.html", map[string]interface{}{
		"emailValue":      emailValue,
		"loginFormMethod": routeMap.Get(entity.PostLoginRouteName, "method"),
		"loginFormAction": routeMap.Get(entity.PostLoginRouteName, "path"),
	})
}

func (h *userHandler) PostLogin(e echo.Context) error {
	// init session
	sess, err := helper.NewSession(e.Request(), e.Response(), h.store, loginSessionName)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  entity.SessionErrorMessage,
			Internal: entity.NewInternalError(entity.SessionError, err.Error()),
		}
	}

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	loginPathRedirect := routeMap.Get(entity.LoginRouteName, "path")
	profilePathRedirect := routeMap.Get(entity.ProfileRouteName, "path")

	// bind payload
	payload := new(entity.AttemptLoginPayload)
	err = e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	// get user
	user, err := h.userModule.AttemptLogin(ctx, payload)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("erorr", herr.Message)
			return e.Redirect(http.StatusMovedPermanently, loginPathRedirect)
		}
		return err
	}

	// generate token
	token, err := h.userModule.GenerateJwtToken(ctx, user)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("erorr", herr.Message)
			return e.Redirect(http.StatusMovedPermanently, loginPathRedirect)
		}
		return err
	}

	// init token session
	tokenSess, err := helper.NewSession(e.Request(), e.Response(), h.store, entity.SessionID)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  entity.SessionErrorMessage,
			Internal: entity.NewInternalError(entity.SessionError, err.Error()),
		}
	}
	tokenSess.Set(entity.TokenSessionName, token)

	return e.Redirect(http.StatusMovedPermanently, profilePathRedirect)
}

func (h *userHandler) PostLogout(e echo.Context) error {
	return e.NoContent(http.StatusNotFound)
}

func (h *userHandler) Profile(e echo.Context) error {
	return e.JSON(http.StatusOK, map[string]interface{}{"user": e.Get(h.cfg.JwtUserContextKey)})
}
