package api

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UserAPI interface {
	// create user
	Create(e echo.Context) error
	// post login
	Login(e echo.Context) error
}

func NewUserHandler(cfg config.Config, userUsecase contract.UserUsecase, mailUsecase contract.MailUsecase) UserAPI {
	return &userHandler{cfg, userUsecase, mailUsecase}
}

type userHandler struct {
	cfg         config.Config
	userUsecase contract.UserUsecase
	mailUsecase contract.MailUsecase
}

func (h *userHandler) Create(e echo.Context) error {
	payload := new(entity.CreateUserPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	// get author
	_author := e.Get(h.cfg.JwtUserContextKey)
	author, ok := _author.(*entity.User)
	if !ok {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.GetAuthorFromJwtError, "failed type assertion author"),
		}
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
