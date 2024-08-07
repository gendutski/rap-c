package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type AuthPage interface {
	// login page
	Login(e echo.Context) error
	// set token to session
	SubmitToken(e echo.Context) error
	// post logout
	Logout(e echo.Context) error
	// password must change page
	PasswordMustChange(e echo.Context) error
	// forgot password page
	ForgotPassword(e echo.Context) error
	// reset password page
	ResetPassword(e echo.Context) error
}

func NewAuthPage(cfg *config.Config, router *config.Route, authUsecase contract.AuthUsecase, sessionUsecase contract.SessionUsecase, mailUsecase contract.MailUsecase) AuthPage {
	return &authHandler{
		cfg:            cfg,
		router:         router,
		authUsecase:    authUsecase,
		sessionUsecase: sessionUsecase,
		mailUsecase:    mailUsecase,
		BaseHandler:    handler.NewBaseHandler(cfg, router),
	}
}

type authHandler struct {
	cfg            *config.Config
	router         *config.Route
	authUsecase    contract.AuthUsecase
	sessionUsecase contract.SessionUsecase
	mailUsecase    contract.MailUsecase
	BaseHandler    *handler.BaseHandler
}

func (h *authHandler) Login(e echo.Context) error {
	// check if token session is exists and valid
	_, _, err := h.sessionUsecase.ValidateJwtToken(e, h.cfg.EnableGuestLogin())
	if err == nil {
		return e.Redirect(http.StatusFound, h.router.DefaultAuthorizedWebPage(
			h.sessionUsecase.GetPrevRoute(e),
		).Path())
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
	// get error stored in session
	herr := h.sessionUsecase.GetError(e)
	if herr != nil {
		if herr.Code != http.StatusUnauthorized {
			return herr
		}
		infos = append(infos, "Sesi telah berakhir, silahkan login ulang!")
	}
	// set infos
	infoMessages := []byte("[]")
	if len(infos) > 0 {
		infoMessages, _ = json.Marshal(infos)
	}

	return e.Render(http.StatusOK, "login.html", map[string]interface{}{
		"enableGuest":              h.cfg.EnableGuestLogin,
		"emailValue":               emailValue,
		"passwordValue":            passValue,
		"infoMessages":             string(infoMessages),
		"guestLoginMethod":         h.router.GuestLoginAPI.Method(),
		"guestLoginAction":         h.cfg.URL(h.router.GuestLoginAPI.Path()),
		"loginFormMethod":          h.router.LoginAPI.Method(),
		"loginFormAction":          h.cfg.URL(h.router.LoginAPI.Path()),
		"submitTokenSessionMethod": h.router.SubmitTokenSessionWebPage.Method(),
		"submitTokenSessionAction": h.cfg.URL(h.router.SubmitTokenSessionWebPage.Path()),
		"forgotPasswordPath":       h.router.ForgotPasswordWebPage.Path(),
	})
}

func (h *authHandler) SubmitToken(e echo.Context) error {
	token := e.FormValue("token")
	if token == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "token empty",
		}
	}
	err := h.sessionUsecase.SaveJwtToken(e, token)
	if err != nil {
		return err
	}

	return e.Redirect(http.StatusFound, h.router.DefaultAuthorizedWebPage(
		h.sessionUsecase.GetPrevRoute(e),
	).Path())
}

func (h *authHandler) Logout(e echo.Context) error {
	// destroy session
	err := h.sessionUsecase.Logout(e)
	if err != nil {
		return err
	}

	// redirect
	return e.Redirect(http.StatusFound, h.router.LoginWebPage.Path())
}

func (h *authHandler) PasswordMustChange(e echo.Context) error {
	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}
	// get token
	token, err := h.BaseHandler.GetToken(e)
	if err != nil {
		return err
	}

	return e.Render(http.StatusOK, "pass-changer.html", map[string]interface{}{
		"emailValue":   author.Email,
		"token":        token,
		"renewAction":  h.cfg.URL(h.router.PasswordMustChangeAPI.Path()),
		"renewMethod":  h.router.PasswordMustChangeAPI.Method(),
		"logoutAction": h.cfg.URL(h.router.LogoutWebPage.Path()),
		"logoutMethod": h.router.LogoutWebPage.Method(),
		"redirectPath": h.router.DefaultAuthorizedWebPage(
			h.sessionUsecase.GetPrevRoute(e),
		).Path(),
	})
}

func (h *authHandler) ForgotPassword(e echo.Context) error {
	// check if token session is exists and valid
	_, _, err := h.sessionUsecase.ValidateJwtToken(e, h.cfg.EnableGuestLogin())
	if err == nil {
		return e.Redirect(http.StatusFound, h.router.DefaultAuthorizedWebPage(
			h.sessionUsecase.GetPrevRoute(e),
		).Path())
	}

	return e.Render(http.StatusOK, "forgot-password.html", map[string]interface{}{
		"formMethod": h.router.RequestResetPasswordAPI.Method(),
		"formAction": h.router.RequestResetPasswordAPI.Path(),
		"loginPath":  h.router.LoginWebPage.Path(),
	})
}

func (h *authHandler) ResetPassword(e echo.Context) error {
	// check if token session is exists and valid
	_, _, err := h.sessionUsecase.ValidateJwtToken(e, h.cfg.EnableGuestLogin())
	if err == nil {
		return e.Redirect(http.StatusFound, h.router.DefaultAuthorizedWebPage(
			h.sessionUsecase.GetPrevRoute(e),
		).Path())
	}

	// bind payload
	payload := new(payloadentity.ValidateResetTokenPayload)
	err = e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError,
				fmt.Sprintf("auth-page.ResetPassword bind error: %v", err)),
		}
	}

	// validate token
	ctx := e.Request().Context()
	err = h.authUsecase.ValidateResetToken(ctx, payload)
	if err != nil {
		return err
	}

	return e.Render(http.StatusOK, "reset-password.html", map[string]interface{}{
		"email":                    payload.Email,
		"token":                    payload.Token,
		"passwordMethod":           h.router.ResetPasswordAPI.Method(),
		"passwordAction":           h.router.ResetPasswordAPI.Path(),
		"submitTokenSessionMethod": h.router.SubmitTokenSessionWebPage.Method(),
		"submitTokenSessionAction": h.cfg.URL(h.router.SubmitTokenSessionWebPage.Path()),
	})
}
