package entity

import (
	"net/http"
	"rap-c/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func InitSession(r *http.Request, w http.ResponseWriter, store sessions.Store, sessionName string, logMode config.LogMode, enableWarnFileLog bool) *Session {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		InitLog(r.RequestURI, r.Method, "get session", http.StatusUnauthorized, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  SessionErrorMessage,
			Internal: NewInternalError(SessionError, err.Error()),
		}, logMode, enableWarnFileLog).Log()
		sess, _ = store.New(r, sessionName)
	}

	return &Session{
		store: store,
		sess:  sess,
		r:     r,
		w:     w,
	}
}

type Session struct {
	store sessions.Store
	sess  *sessions.Session
	r     *http.Request
	w     http.ResponseWriter
}

func (s *Session) Set(key string, value interface{}) {
	s.sess.Values[key] = value
	s.sess.Save(s.r, s.w)
}

func (s *Session) Get(key string) interface{} {
	if len(s.sess.Values) == 0 {
		return nil
	}
	result, ok := s.sess.Values[key]
	if !ok {
		return nil
	}
	return result
}

func (s *Session) Flash(key string) interface{} {
	if len(s.sess.Values) == 0 {
		return nil
	}
	result, ok := s.sess.Values[key]
	if !ok {
		return nil
	}
	s.Remove(key)
	return result
}

func (s *Session) Remove(key string) {
	delete(s.sess.Values, key)
	s.sess.Save(s.r, s.w)
}

func (s *Session) Destroy() {
	s.sess.Options.MaxAge = -1
	s.sess.Save(s.r, s.w)
	s = nil
}
