package sessionusecase

import (
	"net/http"
	"rap-c/app/entity"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	sessionID string = "SESSION_ID"
	tokenKey  string = "token"
	errorKey  string = "error"
	infoKey   string = "info"
)

func (uc *usecase) initSession(r *http.Request) *sessions.Session {
	sess, err := uc.store.Get(r, sessionID)
	if err != nil {
		entity.InitLog(r.RequestURI, r.Method, "get session", http.StatusUnauthorized, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, uc.cfg.LogMode(), uc.cfg.EnableWarnFileLog()).Log()
		sess, _ = uc.store.New(r, sessionID)
	}
	return sess
}
