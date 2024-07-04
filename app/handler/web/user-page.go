package web

import (
	"errors"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/usecase/contract"
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
	// password must change page
	PasswordChanger(e echo.Context) error
	// profile page
	Profile(e echo.Context) error
}

func NewUserPage(cfg config.Config, store sessions.Store, userUsecase contract.UserUsecase) UserPage {
	return &userHandler{cfg, store, userUsecase}
}

type userHandler struct {
	cfg         config.Config
	store       sessions.Store
	userUsecase contract.UserUsecase
}

func (h *userHandler) Login(e echo.Context) error {
	// init router
	routeMap := helper.RouteMap(e.Echo().Routes())
	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	// check if token session is exists
	ctx := e.Request().Context()
	_, err := h.userUsecase.ValidateSessionJwtToken(ctx, e.Request(), e.Response(), h.store, h.cfg.EnableGuestLogin)
	if err == nil {
		return e.Redirect(http.StatusMovedPermanently, authorizedPathRedirect)
	}

	// load session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.EnableWarnFileLog)

	// get email has been inputed in query params or from prev login page
	emailValue := e.QueryParam("email")
	if emailValue == "" {
		if _val, ok := sess.Flash("email").(string); ok {
			emailValue = _val
		}
	}
	// get password inputed in query params
	passValue := e.QueryParam("password")

	return e.Render(http.StatusOK, "login.html", map[string]interface{}{
		"emailValue":      emailValue,
		"passwordValue":   passValue,
		"loginFormMethod": routeMap.Get(entity.PostLoginRouteName, "method"),
		"loginFormAction": routeMap.Get(entity.PostLoginRouteName, "path"),
	})
}

func (h *userHandler) PostLogin(e echo.Context) error {
	// init session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.EnableWarnFileLog)

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	loginPathRedirect := routeMap.Get(entity.LoginRouteName, "path")
	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	// bind payload
	payload := new(entity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	// get user
	user, err := h.userUsecase.AttemptLogin(ctx, payload)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("erorr", herr.Message)
			return e.Redirect(http.StatusMovedPermanently, loginPathRedirect)
		}
		return err
	}

	// generate token
	token, err := h.userUsecase.GenerateJwtToken(ctx, user)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("erorr", herr.Message)
			return e.Redirect(http.StatusMovedPermanently, loginPathRedirect)
		}
		return err
	}

	// init token session
	tokenSess := entity.InitSession(e.Request(), e.Response(), h.store, entity.SessionID, h.cfg.EnableWarnFileLog)
	tokenSess.Set(entity.TokenSessionName, token)

	return e.Redirect(http.StatusMovedPermanently, authorizedPathRedirect)
}

func (h *userHandler) PostLogout(e echo.Context) error {
	// init session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, entity.SessionID, h.cfg.EnableWarnFileLog)
	sess.Destroy()

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	loginPathRedirect := routeMap.Get(entity.LoginRouteName, "path")

	// redirect
	return e.Redirect(http.StatusMovedPermanently, loginPathRedirect)
}

func (h *userHandler) PasswordChanger(e echo.Context) error {
	user, ok := e.Get(h.cfg.JwtUserContextKey).(*entity.User)
	if !ok {
		return errors.New("invalid user")
	}
	return e.Render(http.StatusOK, "pass-changer.html", map[string]interface{}{"emailValue": user.Email})
}

func (h *userHandler) Profile(e echo.Context) error {
	user, ok := e.Get(h.cfg.JwtUserContextKey).(*entity.User)
	if !ok {
		return errors.New("invalid user")
	}

	return e.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"user": user,
	})
}
