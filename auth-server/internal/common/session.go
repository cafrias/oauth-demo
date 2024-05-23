package common

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var sessionKey = "session"

var defaultSessionOptions = &sessions.Options{
	Path:     "/",
	MaxAge:   86400,
	HttpOnly: true,
	Secure:   true,
	SameSite: http.SameSiteLaxMode,
}

// TODO: move secret to a configuration file
var store = sessions.NewCookieStore([]byte("secret"))

type Session interface {
	GetUserInfo() UserInfo
	SetUserInfo(user UserInfo)
	Save(r *http.Request, w http.ResponseWriter) error
	Delete(r *http.Request, w http.ResponseWriter) error
}

type defaultSession struct {
	session *sessions.Session
}

func (s *defaultSession) GetUserInfo() UserInfo {
	id, ok := s.session.Values["user_id"]
	if !ok {
		id = ""
	}

	email, ok := s.session.Values["email"]
	if !ok {
		email = ""
	}

	return UserInfo{
		UserID: id.(string),
		Email:  email.(string),
	}
}

func (s *defaultSession) SetUserInfo(user UserInfo) {
	s.session.Values["user_id"] = user.UserID
	s.session.Values["email"] = user.Email
}

func (s *defaultSession) Save(r *http.Request, w http.ResponseWriter) error {
	s.session.Options = defaultSessionOptions
	return s.session.Save(r, w)
}

func (s *defaultSession) Delete(r *http.Request, w http.ResponseWriter) error {
	s.session.Options.MaxAge = -1
	return s.session.Save(r, w)
}
