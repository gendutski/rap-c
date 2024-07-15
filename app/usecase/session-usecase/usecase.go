package sessionusecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func NewUsecase(cfg *config.Config, store sessions.Store, authUsecase contract.AuthUsecase) contract.SessionUsecase {
	return &usecase{cfg, store, authUsecase}
}

type usecase struct {
	cfg         *config.Config
	store       sessions.Store
	authUsecase contract.AuthUsecase
}

func (uc *usecase) SaveJwtToken(e echo.Context, token string) error {
	r := e.Request()
	sess := uc.initSession(r)
	sess.Values[tokenKey] = token
	err := sess.Save(r, e.Response())
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseSaveJwtTokenError, err.Error()),
		}
	}
	return nil
}

func (uc *usecase) ValidateJwtToken(e echo.Context, guestAccepted bool) (*databaseentity.User, string, error) {
	// get token string from session
	sess := uc.initSession(e.Request())
	sessionToken, ok := sess.Values[tokenKey]
	if !ok {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.SessionTokenNotFoundMessage,
			Internal: entity.NewInternalError(entity.SessionTokenNotFound, entity.SessionTokenNotFoundMessage),
		}
	}
	tokenStr, ok := sessionToken.(string)
	if !ok {
		return nil, "", &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseTokenInvalidType,
				fmt.Sprintf("session token conversion is %T, not string", sessionToken)),
		}
	}

	// parse token string
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.JwtSecret()), nil
	})
	if err != nil || !token.Valid {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.ValidateTokenFailedMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenFailed, err.Error()),
		}
	}

	// validate token
	user, err := uc.authUsecase.ValidateJwtToken(e.Request().Context(), token, guestAccepted)
	if err != nil {
		return nil, "", err
	}
	return user, tokenStr, nil
}

func (uc *usecase) SetError(e echo.Context, theError error) error {
	// get status code
	herr, ok := theError.(*echo.HTTPError)
	if !ok {
		herr = echo.NewHTTPError(http.StatusInternalServerError)
	}

	// set error log
	entity.InitLog(
		e.Request().RequestURI,
		e.Request().Method,
		"session",
		herr.Code,
		herr,
		uc.cfg.LogMode(),
		uc.cfg.EnableWarnFileLog(),
	).Log()

	// get or set session
	sess := uc.initSession(e.Request())
	internal, ok := herr.Internal.(*entity.InternalError)
	if !ok {
		internal = &entity.InternalError{}
	}
	sess.Values[errorKey], _ = json.Marshal(sessionError{
		Code:     herr.Code,
		Message:  herr.Message,
		Internal: internal,
	})
	err := sess.Save(e.Request(), e.Response())
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseSetErrorError, err.Error()),
		}
	}
	return nil
}

func (uc *usecase) GetError(e echo.Context) *echo.HTTPError {
	// get or set session
	sess := uc.initSession(e.Request())
	sessionValues, ok := sess.Values[errorKey]
	if !ok {
		// no error sent
		return nil
	}
	// get session value
	strJSON, ok := sessionValues.([]byte)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseErrorInvalidType,
				fmt.Sprintf("session error conversion is %T, not []byte", sessionValues)),
		}
	}
	// unmarshal
	var herr sessionError
	err := json.Unmarshal(strJSON, &herr)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseGetErrorError, err.Error()),
		}
	}
	// delete error from session
	delete(sess.Values, errorKey)
	sess.Save(e.Request(), e.Response())

	// return
	return &echo.HTTPError{
		Code:     herr.Code,
		Message:  herr.Message,
		Internal: herr.Internal,
	}
}

func (uc *usecase) SetInfo(e echo.Context, info interface{}) error {
	r := e.Request()
	sess := uc.initSession(r)
	sess.Values[infoKey] = info
	err := sess.Save(r, e.Response())
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseSetInfoError, err.Error()),
		}
	}
	return nil
}

func (uc *usecase) GetInfo(e echo.Context) (interface{}, error) {
	r := e.Request()
	sess := uc.initSession(r)
	info, ok := sess.Values[infoKey]
	if !ok {
		return nil, nil
	}

	// delete session
	delete(sess.Values, infoKey)
	err := sess.Save(r, e.Response())
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseGetInfoError, err.Error()),
		}
	}
	return info, nil
}

func (uc *usecase) SetPrevRoute(e echo.Context) error {
	r := e.Request()
	sess := uc.initSession(r)
	sess.Values[prevRouteKey], _ = json.Marshal(map[string]string{
		prevRouteMapMethod: r.Method,
		prevRouteMapPath:   r.RequestURI,
	})
	err := sess.Save(r, e.Response())
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseSetPrevRouteError, err.Error()),
		}
	}
	return nil
}

func (uc *usecase) GetPrevRoute(e echo.Context) (method string, path string) {
	r := e.Request()

	// get session, if not exists go away
	sess := uc.initSession(r)
	prev, ok := sess.Values[prevRouteKey]
	if !ok {
		return
	}
	// check session assertion
	strJSON, ok := prev.([]byte)
	if !ok {
		entity.InitLog(r.RequestURI, r.Method, "get session", http.StatusInternalServerError, &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseGetInfoError,
				fmt.Sprintf("session error conversion is %T, not []byte", prev)),
		}, uc.cfg.LogMode(), uc.cfg.EnableWarnFileLog()).Log()
		return
	}

	// delete session
	delete(sess.Values, prevRouteKey)
	err := sess.Save(r, e.Response())
	if err != nil {
		entity.InitLog(r.RequestURI, r.Method, "get session", http.StatusInternalServerError, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseGetInfoError, err.Error()),
		}, uc.cfg.LogMode(), uc.cfg.EnableWarnFileLog()).Log()
		return
	}

	// unmarshal prev
	var result map[string]string
	err = json.Unmarshal(strJSON, &result)
	if err != nil {
		entity.InitLog(r.RequestURI, r.Method, "get session", http.StatusInternalServerError, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseGetInfoError, err.Error()),
		}, uc.cfg.LogMode(), uc.cfg.EnableWarnFileLog()).Log()
		return
	}
	method = result[prevRouteMapMethod]
	path = result[prevRouteMapPath]
	return
}

func (uc *usecase) Logout(e echo.Context) error {
	sess := uc.initSession(e.Request())
	sess.Options.MaxAge = -1
	err := sess.Save(e.Request(), e.Response())
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.SessionUsecaseLogoutError, err.Error()),
		}
	}
	return nil
}
