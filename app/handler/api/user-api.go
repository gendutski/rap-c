package api

import (
	"net/http"
	"rap-c/app/entity"
	usermodule "rap-c/app/module/user-module"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UserAPI interface {
	// create user
	Create(e echo.Context) error
	// post login
	Login(e echo.Context) error
}

func NewUserHandler(cfg config.Config, userModule usermodule.UserUsecase) UserAPI {
	return &userHandler{cfg, userModule}
}

type userHandler struct {
	cfg        config.Config
	userModule usermodule.UserUsecase
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
	user, err := h.userModule.Create(ctx, payload, author)
	if err != nil {
		return err
	}

	// send email
	// TODO:

	return e.JSON(http.StatusOK, user)
}

func (h *userHandler) Login(e echo.Context) error {
	payload := new(entity.AttemptLoginPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	ctx := e.Request().Context()

	user, err := h.userModule.AttemptLogin(ctx, payload)
	if err != nil {
		return err
	}

	token, err := h.userModule.GenerateJwtToken(ctx, user)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
