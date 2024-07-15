package web

import (
	"encoding/json"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type AuthPage interface {
	// login page
	Login(e echo.Context) error
	// // post logout
	// PostLogout(e echo.Context) error
	// // password must change page
	// PasswordChanger(e echo.Context) error
	// // request reset password page
	// RequestResetPassword(e echo.Context) error
	// // reset password page
	// ResetPassword(e echo.Context) error
}

func NewAuthPage(cfg *config.Config, store sessions.Store, authUsecase contract.AuthUsecase, sessionUsecase contract.SessionUsecase, mailUsecase contract.MailUsecase) AuthPage {
	return &authHandler{
		cfg:            cfg,
		store:          store,
		authUsecase:    authUsecase,
		sessionUsecase: sessionUsecase,
		mailUsecase:    mailUsecase,
		BaseHandler:    handler.NewBaseHandler(cfg),
	}
}

type authHandler struct {
	cfg            *config.Config
	store          sessions.Store
	authUsecase    contract.AuthUsecase
	sessionUsecase contract.SessionUsecase
	mailUsecase    contract.MailUsecase
	BaseHandler    *handler.BaseHandler
}

func (h *authHandler) Login(e echo.Context) error {
	// check if token session is exists and valid
	_, _, err := h.sessionUsecase.ValidateJwtToken(e, h.cfg.EnableGuestLogin())
	if err == nil {
		return e.Redirect(http.StatusFound, entity.WebDefaultAuthorizedPath)
	}

	// get email has been inputed in query params
	emailValue := e.QueryParam("email")
	// get password inputed in query params
	passValue := e.QueryParam("password")

	// get info
	var infos []string
	sessInfo, err := h.sessionUsecase.GetInfo(e)
	if err != nil {
		return err
	}
	switch inf := sessInfo.(type) {
	case string:
		infos = append(infos, inf)
	case []string:
		infos = append(infos, inf...)
	}
	infoMessages := []byte("[]")
	if len(infos) > 0 {
		infoMessages, _ = json.Marshal(infos)
	}

	// get error stored in session
	herr := h.sessionUsecase.GetError(e)
	if herr != nil {
		if herr.Code != http.StatusUnauthorized {
			return herr
		}
		infos = append(infos, "Sesi telah berakhir, silahkan login ulang!")
	}

	return e.Render(http.StatusOK, "login.html", map[string]interface{}{
		"enableGuest":   h.cfg.EnableGuestLogin,
		"emailValue":    emailValue,
		"passwordValue": passValue,
		"infoMessages":  string(infoMessages),
		// "guestLoginMethod": routeMap.Get(entity.ApiGuestLoginRouteName, "method"),
		// "guestLoginAction": routeMap.Get(entity.ApiGuestLoginRouteName, "path"),
		// "loginFormMethod":  routeMap.Get(entity.ApiLoginRouteName, "method"),
		// "loginFormAction":  routeMap.Get(entity.ApiLoginRouteName, "path"),
		"resetPath": entity.WebResetPasswordPath,
	})
}

// func (h *authHandler) PostLogout(e echo.Context) error {
// 	// destroy session
// 	err := h.sessionUsecase.Logout(e)
// 	if err != nil {
// 		return err
// 	}

// 	// redirect
// 	return e.Redirect(http.StatusFound, entity.WebLoginPath)
// }

// func (h *authHandler) PasswordChanger(e echo.Context) error {
// 	// get author
// 	author, err := h.BaseHandler.GetAuthor(e)
// 	if err != nil {
// 		return err
// 	}
// 	// get token
// 	token, err := h.BaseHandler.GetToken(e)
// 	if err != nil {
// 		return err
// 	}

// 	return e.Render(http.StatusOK, "pass-changer.html", map[string]interface{}{
// 		"emailValue":   author.Email,
// 		"token":        token,
// 		"logoutAction": entity.WebLogoutPath,
// 	})
// }

// func (h *authHandler) RequestResetPassword(e echo.Context) error {
// 	// init router
// 	routeMap := helper.RouteMap(e.Echo().Routes())
// 	authorizedPathRedirect := routeMap.Get(entity.DefaultAuthorizedRouteRedirect, "path")

// 	// check if token session is exists
// 	ctx := e.Request().Context()
// 	_, _, err := h.authUsecase.ValidateSessionJwtToken(ctx, e.Request(), e.Response(), h.store, h.cfg.EnableGuestLogin)
// 	if err == nil {
// 		return e.Redirect(http.StatusFound, authorizedPathRedirect)
// 	}

// 	// load session
// 	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.LogMode, h.cfg.EnableWarnFileLog)

// 	// get submit login error
// 	var submitErr []string
// 	var submitError []byte = []byte("[]")
// 	if _submitErr := sess.Flash("error"); _submitErr != nil {
// 		switch _val := _submitErr.(type) {
// 		case []string:
// 			submitErr = append(submitErr, _val...)
// 		case string:
// 			submitErr = append(submitErr, _val)
// 		}
// 	}
// 	if len(submitErr) > 0 {
// 		submitError, _ = json.Marshal(submitErr)
// 	}

// 	return e.Render(http.StatusOK, "request-reset.html", map[string]interface{}{
// 		"enableGuest": h.cfg.EnableGuestLogin,
// 		"loginError":  string(submitError),
// 		"formMethod":  routeMap.Get(entity.PostRequestResetPasswordName, "method"),
// 		"formAction":  routeMap.Get(entity.PostRequestResetPasswordName, "path"),
// 		"loginPath":   routeMap.Get(entity.LoginRouteName, "path"),
// 	})
// }

// func (h *authHandler) ResetPassword(e echo.Context) error {
// 	ctx := e.Request().Context()
// 	email := e.QueryParam("email")
// 	token := e.QueryParam("token")

// 	err := h.authUsecase.ValidateResetPassword(ctx, email, token)
// 	if err != nil {
// 		return err
// 	}

// 	// map route
// 	routeMap := helper.RouteMap(e.Echo().Routes())

// 	// load session
// 	sess := entity.InitSession(e.Request(), e.Response(), h.store, loginSessionName, h.cfg.LogMode, h.cfg.EnableWarnFileLog)

// 	// get submit login error
// 	var submitErr []string
// 	var submitError []byte = []byte("[]")
// 	if _submitErr := sess.Flash("error"); _submitErr != nil {
// 		switch _val := _submitErr.(type) {
// 		case []string:
// 			submitErr = append(submitErr, _val...)
// 		case string:
// 			submitErr = append(submitErr, _val)
// 		}
// 	}
// 	if len(submitErr) > 0 {
// 		submitError, _ = json.Marshal(submitErr)
// 	}

// 	return e.Render(http.StatusOK, "reset-password.html", map[string]interface{}{
// 		"email":          email,
// 		"token":          token,
// 		"submitError":    string(submitError),
// 		"passwordMethod": routeMap.Get(entity.SubmitResetPasswordName, "method"),
// 		"passwordAction": routeMap.Get(entity.SubmitResetPasswordName, "path"),
// 	})
// }
