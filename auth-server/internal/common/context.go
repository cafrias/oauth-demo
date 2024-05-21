package common

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var sessionKey = "session"

// TODO: move secret to a configuration file
var store = sessions.NewCookieStore([]byte("secret"))

type UserInfo struct {
	UserID string
	Email  string
}

type Session interface {
	GetUserInfo() UserInfo
	SetUserInfo(user UserInfo)
	Save(r *http.Request, w http.ResponseWriter) error
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
	return s.session.Save(r, w)
}

type AppContext interface {
	echo.Context
	GetSession() (Session, error)
	SaveSession(Session) error
}

type defaultAppContext struct {
	echo.Context
}

func (c *defaultAppContext) GetSession() (Session, error) {
	session, err := store.Get(c.Request(), sessionKey)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse session: %w", err)
	}

	return &defaultSession{session}, nil
}

func (c *defaultAppContext) SaveSession(session Session) error {
	return session.Save(c.Request(), c.Response())
}

func UseAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(&defaultAppContext{c})
	}
}
