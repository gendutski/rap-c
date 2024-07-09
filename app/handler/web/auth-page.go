package web

import (
	"encoding/json"
	"log"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler"
	"rap-c/app/helper"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	loginSessionName string = "logsession"
)

type AuthPage interface {
	// login page
	Login(e echo.Context) error
	// guest login
	GuestLogin(e echo.Context) error
	// post login form
	PostLogin(e echo.Context) error
	// post logout
	PostLogout(e echo.Context) error
	// password must change page
	PasswordChanger(e echo.Context) error
	// submit password change
	SubmitPasswordChanger(e echo.Context) error
	// request reset password page
	RequestResetPassword(e echo.Context) error
	// submit request reset password page
	SubmitRequestResetPassword(e echo.Context) error
	// reset password page
	ResetPassword(e echo.Context) error
}

func NewAuthPage(cfg config.Config, store sessions.Store, authUsecase contract.AuthUsecase, mailUsecase contract.MailUsecase) AuthPage {
	return &authHandler{
		cfg:         cfg,
		store:       store,
		authUsecase: authUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg),
	}
}

type authHandler struct {
	cfg         config.Config
	store       sessions.Store
	authUsecase contract.AuthUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *authHandler) Login(e echo.Context) error {
	// init router
	routeMap := helper.RouteMap(e.Echo().Routes())
	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	// check if token session is exists
	ctx := e.Request().Context()
	_, _, err := h.authUsecase.ValidateSessionJwtToken(ctx, e.Request(), e.Response(), h.store, h.cfg.EnableGuestLogin)
	if err == nil {
		return e.Redirect(http.StatusFound, authorizedPathRedirect)
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

	// get submit login error
	var loginErr []string
	var loginError []byte = []byte("[]")
	if _loginErr := sess.Flash("error"); _loginErr != nil {
		switch _val := _loginErr.(type) {
		case []string:
			loginErr = append(loginErr, _val...)
		case string:
			loginErr = append(loginErr, _val)
		}
	}
	if len(loginErr) > 0 {
		loginError, _ = json.Marshal(loginErr)
	}

	// get info
	var infoMsg []string
	var infoMessages []byte = []byte("[]")
	if _infoMsg := sess.Flash("info"); _infoMsg != nil {
		switch _val := _infoMsg.(type) {
		case []string:
			infoMsg = append(infoMsg, _val...)
		case string:
			infoMsg = append(infoMsg, _val)
		}
	}
	if len(infoMsg) > 0 {
		infoMessages, _ = json.Marshal(infoMsg)
	}

	return e.Render(http.StatusOK, "login.html", map[string]interface{}{
		"enableGuest":      h.cfg.EnableGuestLogin,
		"emailValue":       emailValue,
		"passwordValue":    passValue,
		"loginError":       string(loginError),
		"infoMessages":     string(infoMessages),
		"guestLoginMethod": routeMap.Get(entity.GuestLoginRouteName, "method"),
		"guestLoginAction": routeMap.Get(entity.GuestLoginRouteName, "path"),
		"loginFormMethod":  routeMap.Get(entity.PostLoginRouteName, "method"),
		"loginFormAction":  routeMap.Get(entity.PostLoginRouteName, "path"),
		"resetPath":        routeMap.Get(entity.RequestResetPasswordName, "path"),
	})
}

func (h *authHandler) GuestLogin(e echo.Context) error {
	ctx := e.Request().Context()

	// get user
	user, err := h.authUsecase.AttemptGuestLogin(ctx)
	if err != nil {
		return err
	}

	// generate token
	token, err := h.authUsecase.GenerateJwtToken(ctx, user, false)
	if err != nil {
		return err
	}

	// init token session
	tokenSess := entity.InitSession(e.Request(), e.Response(), h.store, entity.SessionID, h.cfg.EnableWarnFileLog)
	tokenSess.Set(entity.TokenSessionName, token)

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	return e.Redirect(http.StatusFound, authorizedPathRedirect)
}

func (h *authHandler) PostLogin(e echo.Context) error {
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
	user, err := h.authUsecase.AttemptLogin(ctx, payload)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("error", herr.Message)
			return e.Redirect(http.StatusFound, loginPathRedirect)
		}
		return err
	}

	// generate token
	token, err := h.authUsecase.GenerateJwtToken(ctx, user, payload.RememberMe)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("email", e.FormValue("email"))
			sess.Set("error", herr.Message)
			return e.Redirect(http.StatusFound, loginPathRedirect)
		}
		return err
	}

	// init token session
	tokenSess := entity.InitSession(e.Request(), e.Response(), h.store, entity.SessionID, h.cfg.EnableWarnFileLog)
	tokenSess.Set(entity.TokenSessionName, token)

	return e.Redirect(http.StatusFound, authorizedPathRedirect)
}

func (h *authHandler) PostLogout(e echo.Context) error {
	// init session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, entity.SessionID, h.cfg.EnableWarnFileLog)
	sess.Destroy()

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	loginPathRedirect := routeMap.Get(entity.LoginRouteName, "path")

	// redirect
	return e.Redirect(http.StatusFound, loginPathRedirect)
}

func (h *authHandler) PasswordChanger(e echo.Context) error {
	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())

	return e.Render(http.StatusOK, "pass-changer.html", map[string]interface{}{
		"emailValue":   author.Email,
		"logoutAction": routeMap.Get(entity.PostLogoutRouteName, "path"),
		"logoutMethod": routeMap.Get(entity.PostLogoutRouteName, "method"),
		"renewAction":  routeMap.Get(entity.RenewPasswordRouteName, "path"),
		"renewMethod":  routeMap.Get(entity.RenewPasswordRouteName, "method"),
	})
}

func (h *authHandler) SubmitPasswordChanger(e echo.Context) error {
	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}

	// get payload
	payload := new(entity.RenewPasswordPayload)
	err = e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	// renew password
	err = h.authUsecase.RenewPassword(ctx, author, payload)
	if err != nil {
		return err
	}

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	defaultRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	// redirect
	return e.Redirect(http.StatusFound, defaultRedirect)
}

func (h *authHandler) RequestResetPassword(e echo.Context) error {
	// init router
	routeMap := helper.RouteMap(e.Echo().Routes())
	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

	// check if token session is exists
	ctx := e.Request().Context()
	_, _, err := h.authUsecase.ValidateSessionJwtToken(ctx, e.Request(), e.Response(), h.store, h.cfg.EnableGuestLogin)
	if err == nil {
		return e.Redirect(http.StatusFound, authorizedPathRedirect)
	}

	// load session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.EnableWarnFileLog)

	// get submit login error
	var submitErr []string
	var submitError []byte = []byte("[]")
	if _submitErr := sess.Flash("error"); _submitErr != nil {
		switch _val := _submitErr.(type) {
		case []string:
			submitErr = append(submitErr, _val...)
		case string:
			submitErr = append(submitErr, _val)
		}
	}
	if len(submitErr) > 0 {
		submitError, _ = json.Marshal(submitErr)
	}

	return e.Render(http.StatusOK, "request-reset.html", map[string]interface{}{
		"enableGuest": h.cfg.EnableGuestLogin,
		"loginError":  string(submitError),
		"formMethod":  routeMap.Get(entity.PostRequestResetPasswordName, "method"),
		"formAction":  routeMap.Get(entity.PostRequestResetPasswordName, "path"),
		"loginPath":   routeMap.Get(entity.LoginRouteName, "path"),
	})
}

func (h *authHandler) SubmitRequestResetPassword(e echo.Context) error {
	// init session
	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.EnableWarnFileLog)

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())
	loginPathRedirect := routeMap.Get(entity.LoginRouteName, "path")
	requestResetRedirect := routeMap.Get(entity.RequestResetPasswordName, "path")

	// bind payload
	payload := new(entity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	// get user & token
	user, token, err := h.authUsecase.RequestResetPassword(ctx, e.FormValue("email"))
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("error", herr.Message)
			return e.Redirect(http.StatusFound, requestResetRedirect)
		}
		return err
	}

	// send email
	err = h.mailUsecase.ResetPassword(user, token)
	if err != nil {
		if herr, ok := err.(*echo.HTTPError); ok {
			sess.Set("error", herr.Message)
			return e.Redirect(http.StatusFound, requestResetRedirect)
		}
		return err
	}

	sess.Set("info", "email for request reset password has been sent")
	return e.Redirect(http.StatusFound, loginPathRedirect)
}

func (h *authHandler) ResetPassword(e echo.Context) error {
	ctx := e.Request().Context()
	email := e.QueryParam("email")
	token := e.QueryParam("token")
	log.Println("FUCK")

	err := h.authUsecase.ValidateResetPassword(ctx, email, token)
	if err != nil {
		return err
	}

	// map route
	routeMap := helper.RouteMap(e.Echo().Routes())

	return e.Render(http.StatusOK, "reset-password.html", map[string]interface{}{
		"email":          email,
		"token":          token,
		"passwordMethod": routeMap.Get(entity.SubmitResetPasswordName, "method"),
		"passwordAction": routeMap.Get(entity.SubmitResetPasswordName, "path"),
	})
}
