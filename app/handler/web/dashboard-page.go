package web

import (
	"net/http"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type DashboardPage interface {
	// profile page
	Dashboard(e echo.Context) error
}

func NewDashboardPage(cfg *config.Config, router *config.Route, sessionUsecase contract.SessionUsecase) DashboardPage {
	return &dashboardHandler{
		cfg:            cfg,
		router:         router,
		sessionUsecase: sessionUsecase,
		BaseHandler:    handler.NewBaseHandler(cfg, router),
	}
}

type dashboardHandler struct {
	cfg            *config.Config
	router         *config.Route
	sessionUsecase contract.SessionUsecase
	BaseHandler    *handler.BaseHandler
}

func (h *dashboardHandler) Dashboard(e echo.Context) error {
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

	return e.Render(http.StatusOK, "dashboard", map[string]interface{}{
		"author":  author,
		"token":   token,
		"title":   "Dashboard",
		"layouts": h.BaseHandler.GetLayouts("profile"),
	})
}
