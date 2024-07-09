package web

import (
	"net/http"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type UserPage interface {
	// profile page
	Profile(e echo.Context) error
}

func NewUserPage(cfg config.Config, store sessions.Store, userUsecase contract.UserUsecase, mailUsecase contract.MailUsecase) UserPage {
	return &userHandler{
		cfg:         cfg,
		store:       store,
		userUsecase: userUsecase,
		mailUsecase: mailUsecase,
		BaseHandler: handler.NewBaseHandler(cfg),
	}
}

type userHandler struct {
	cfg         config.Config
	store       sessions.Store
	userUsecase contract.UserUsecase
	mailUsecase contract.MailUsecase
	BaseHandler *handler.BaseHandler
}

func (h *userHandler) Profile(e echo.Context) error {
	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}

	return e.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"author":  author,
		"title":   "Profile",
		"layouts": h.BaseHandler.GetLayouts("profile"),
	})
}
