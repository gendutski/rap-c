package helper

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func NewSession(r *http.Request, w http.ResponseWriter, store sessions.Store, sessionName string) (*session, error) {
	sess, err := store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	return &session{
		store: store,
		sess:  sess,
		r:     r,
		w:     w,
	}, nil
}

type session struct {
	store sessions.Store
	sess  *sessions.Session
	r     *http.Request
	w     http.ResponseWriter
}

func (s *session) Set(key string, value interface{}) {
	s.sess.Values[key] = value
	s.sess.Save(s.r, s.w)
}

func (s *session) Get(key string) interface{} {
	if len(s.sess.Values) == 0 {
		return nil
	}
	result, ok := s.sess.Values[key]
	if !ok {
		return nil
	}
	return result
}

func (s *session) Flash(key string) interface{} {
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

func (s *session) Remove(key string) {
	delete(s.sess.Values, key)
	s.sess.Save(s.r, s.w)
}

func (s *session) Destroy() {
	s.sess.Options.MaxAge = -1
	s.sess.Save(s.r, s.w)
	s = nil
}
