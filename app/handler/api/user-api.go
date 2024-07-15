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

type UserAPI interface {
	// create user
	Create(e echo.Context) error
	// get user list
	GetUserList(e echo.Context) error
	// get total user list
	GetTotalUserList(e echo.Context) error
	// get user detail by username
	GetUserDetailByUsername(e echo.Context) error
}

func NewUserHandler(cfg *config.Config, router *config.Route, userUsecase contract.UserUsecase, mailUsecase contract.MailUsecase) UserAPI {
	return &userHandler{
		cfg:         cfg,
		router:      router,
		userUsecase: userUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg, router),
	}
}

type userHandler struct {
	cfg         *config.Config
	router      *config.Route
	userUsecase contract.UserUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *userHandler) Create(e echo.Context) error {
	payload := new(payloadentity.CreateUserPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("user-api.Create bind error: %v", err)),
		}
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
			e.Request().Method,
			fmt.Sprintf("Send welcome email to %s", user.Email),
			http.StatusOK,
			nil,
			h.cfg.LogMode(),
			false,
		).Log()
		err = h.mailUsecase.Welcome(user, password)
		if err != nil {
			entity.InitLog(
				e.Request().RequestURI,
				e.Request().Method,
				"send welcome mail",
				http.StatusOK,
				err,
				h.cfg.LogMode(),
				h.cfg.EnableWarnFileLog(),
			).Log()
		}
	}()

	return e.JSON(http.StatusOK, user)
}

func (h *userHandler) GetUserList(e echo.Context) error {
	req := new(payloadentity.GetUserListRequest)
	err := e.Bind(req)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("user-api.GetUserList bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	users, err := h.userUsecase.GetUserList(ctx, req)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"users":   users,
		"request": req,
	})
}

func (h *userHandler) GetTotalUserList(e echo.Context) error {
	req := new(payloadentity.GetUserListRequest)
	err := e.Bind(req)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("user-api.GetTotalUserList bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	total, err := h.userUsecase.GetTotalUserList(ctx, req)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"total":   total,
		"request": req,
	})
}

func (h *userHandler) GetUserDetailByUsername(e echo.Context) error {
	username := e.Param("username")
	ctx := e.Request().Context()

	user, err := h.userUsecase.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, user)
}
