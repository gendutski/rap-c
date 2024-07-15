package api

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type AuthAPI interface {
	// post login
	Login(e echo.Context) error
	// guest login
	GuestLogin(e echo.Context) error
	// renew password for must change password
	RenewPassword(e echo.Context) error
	// submit request reset/forgot password
	RequestResetPassword(e echo.Context) error
	// do reset password
	ResetPassword(e echo.Context) error
}

func NewAuthHandler(cfg *config.Config, router *config.Route, authUsecase contract.AuthUsecase, mailUsecase contract.MailUsecase) AuthAPI {
	return &authHandler{
		cfg:         cfg,
		router:      router,
		authUsecase: authUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg, router),
	}
}

type authHandler struct {
	cfg         *config.Config
	router      *config.Route
	authUsecase contract.AuthUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *authHandler) Login(e echo.Context) error {
	payload := new(payloadentity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("auth-api.Login bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	user, err := h.authUsecase.AttemptLogin(ctx, payload)
	if err != nil {
		return err
	}

	token, err := h.authUsecase.GenerateJwtToken(ctx, user, payload.RememberMe)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *authHandler) GuestLogin(e echo.Context) error {
	ctx := e.Request().Context()
	user, err := h.authUsecase.AttemptGuestLogin(ctx)
	if err != nil {
		return err
	}

	token, err := h.authUsecase.GenerateJwtToken(ctx, user, false)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *authHandler) RenewPassword(e echo.Context) error {
	payload := new(payloadentity.RenewPasswordPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("auth-api.RenewPassword bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}

	err = h.authUsecase.RenewPassword(ctx, author, payload)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *authHandler) RequestResetPassword(e echo.Context) error {
	payload := new(payloadentity.RequestResetPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("auth-api.RequestResetPassword bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	// get data
	user, token, err := h.authUsecase.RequestResetPassword(ctx, payload)
	if err != nil {
		return err
	}

	// send email
	err = h.mailUsecase.ResetPassword(user, token)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"success": "email for request reset password has been sent",
	})
}

func (h *authHandler) ResetPassword(e echo.Context) error {
	payload := new(payloadentity.ResetPasswordPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("auth-api.ResetPassword bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	// reset password
	user, err := h.authUsecase.SubmitResetPassword(ctx, payload)
	if err != nil {
		return err
	}

	// generate token
	token, err := h.authUsecase.GenerateJwtToken(ctx, user, false)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
