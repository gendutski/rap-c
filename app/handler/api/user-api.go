package api

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UserAPI interface {
	// create user
	Create(e echo.Context) error
	// post login
	Login(e echo.Context) error
	// renew password
	RenewPassword(e echo.Context) error
}

func NewUserHandler(cfg config.Config, userUsecase contract.UserUsecase, mailUsecase contract.MailUsecase) UserAPI {
	return &userHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg),
	}
}

type userHandler struct {
	cfg         config.Config
	userUsecase contract.UserUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *userHandler) Create(e echo.Context) error {
	payload := new(entity.CreateUserPayload)
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

	// create user
	user, password, err := h.userUsecase.Create(ctx, payload, author)
	if err != nil {
		return err
	}

	// send email
	go func() {
		entity.InitLog(
			e.Request().RequestURI,
			fmt.Sprintf("Send welcome email to %s", user.Email),
			http.StatusOK,
			nil,
			false,
		).Log()
		err = h.mailUsecase.Welcome(user, password)
		if err != nil {
			entity.InitLog(
				e.Request().RequestURI,
				"send welcome mail",
				http.StatusOK,
				err,
				h.cfg.EnableWarnFileLog,
			).Log()
		}
	}()

	return e.JSON(http.StatusOK, user)
}

func (h *userHandler) Login(e echo.Context) error {
	payload := new(entity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	user, err := h.userUsecase.AttemptLogin(ctx, payload)
	if err != nil {
		return err
	}

	token, err := h.userUsecase.GenerateJwtToken(ctx, user)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *userHandler) RenewPassword(e echo.Context) error {
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

	err = h.userUsecase.RenewPassword(ctx, author, payload)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}
