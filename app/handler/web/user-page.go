package web

import (
	"net/http"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UserPage interface {
	// profile page
	Profile(e echo.Context) error
}

func NewUserPage(cfg *config.Config, router *config.Route, sessionUsecase contract.SessionUsecase, userUsecase contract.UserUsecase, mailUsecase contract.MailUsecase) UserPage {
	return &userHandler{
		cfg:            cfg,
		router:         router,
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
		mailUsecase:    mailUsecase,
		BaseHandler:    handler.NewBaseHandler(cfg, router),
	}
}

type userHandler struct {
	cfg            *config.Config
	router         *config.Route
	sessionUsecase contract.SessionUsecase
	userUsecase    contract.UserUsecase
	mailUsecase    contract.MailUsecase
	BaseHandler    *handler.BaseHandler
}

func (h *userHandler) Profile(e echo.Context) error {
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

	return e.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"author":           author,
		"token":            token,
		"title":            "Profile",
		"layouts":          h.BaseHandler.GetLayouts("profile"),
		"formUpdateMethod": h.router.UpdateUserAPI.Method(),
		"formUpdateAction": h.router.UpdateUserAPI.Path(),
	})
}
