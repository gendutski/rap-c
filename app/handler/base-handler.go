package handler

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/config"
	"time"

	"github.com/labstack/echo/v4"
)

func NewBaseHandler(cfg config.Config) *BaseHandler {
	return &BaseHandler{cfg}
}

type BaseHandler struct {
	cfg config.Config
}

type Layouts struct {
	AppName      string
	Copyright    string
	LogoutPath   string
	LogoutMethod string
}

func (h *BaseHandler) GetAuthor(e echo.Context) (*entity.User, error) {
	// get author
	_author := e.Get(h.cfg.JwtUserContextKey)
	author, ok := _author.(*entity.User)
	if !ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.GetAuthorFromJwtError, "failed type assertion author"),
		}
	}
	return author, nil
}

func (h *BaseHandler) GetLayouts() Layouts {
	return Layouts{
		AppName:      config.AppName,
		Copyright:    h.GetCopyright(),
		LogoutPath:   entity.WebLogoutPath,
		LogoutMethod: entity.WebLogoutMethod,
	}
}

func (h *BaseHandler) GetCopyright() string {
	return fmt.Sprintf(`Copyright &copy; <a href="https://github.com/gendutski/rap-c" target="_blank">Rap-C</a> %s`, time.Now().Format("2006"))
}
