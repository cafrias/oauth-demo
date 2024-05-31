package common

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type UserInfo struct {
	UserID string
	Email  string
}

type AppContext interface {
	echo.Context
	// GetUser returns the user information stored in the session
	GetUser() (UserInfo, error)
	GetSession() (Session, error)
	SaveSession(Session) error
	DeleteSession(Session) error
}

type defaultAppContext struct {
	echo.Context
}

func (c *defaultAppContext) GetUser() (UserInfo, error) {
	s, err := c.GetSession()
	if err != nil {
		return UserInfo{}, fmt.Errorf("Unable to get session: %w", err)
	}

	return s.GetUserInfo(), nil
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

func (c *defaultAppContext) DeleteSession(session Session) error {
	return session.Delete(c.Request(), c.Response())
}

func UseAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(&defaultAppContext{c})
	}
}
