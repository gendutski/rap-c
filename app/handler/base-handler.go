package handler

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

func NewBaseHandler(cfg config.Config) *BaseHandler {
	return &BaseHandler{cfg}
}

type BaseHandler struct {
	cfg config.Config
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
