package api

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type AuthAPI interface {
	// post login
	Login(e echo.Context) error
	// renew password
	RenewPassword(e echo.Context) error
}

func NewAuthHandler(cfg config.Config, authUsecase contract.AuthUsecase, mailUsecase contract.MailUsecase) AuthAPI {
	return &authHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg),
	}
}

type authHandler struct {
	cfg         config.Config
	authUsecase contract.AuthUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *authHandler) Login(e echo.Context) error {
	payload := new(entity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
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

func (h *authHandler) RenewPassword(e echo.Context) error {
	payload := new(entity.RenewPasswordPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
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
