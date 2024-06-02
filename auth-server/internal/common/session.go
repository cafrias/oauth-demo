package common

import (
	"auth-server/internal/security"
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
	IsAuthenticated() bool
	Save(r *http.Request, w http.ResponseWriter) error
	Delete(r *http.Request, w http.ResponseWriter) error
	GetCSRFToken() string
}

type defaultSession struct {
	session *sessions.Session
}

func NewDefaultSession(r *http.Request) (Session, error) {
	session, err := store.Get(r, sessionKey)
	if err != nil {
		return nil, err
	}

	t, err := security.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}
	session.Values["csrf_token"] = t
	session.Options = defaultSessionOptions

	return &defaultSession{session}, nil
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

func (s *defaultSession) GetCSRFToken() string {
	token, _ := s.session.Values["csrf_token"].(string)
	return token
}

func (s *defaultSession) SetUserInfo(user UserInfo) {
	s.session.Values["user_id"] = user.UserID
	s.session.Values["email"] = user.Email
}

func (s *defaultSession) Save(r *http.Request, w http.ResponseWriter) error {
	return s.session.Save(r, w)
}

func (s *defaultSession) Delete(r *http.Request, w http.ResponseWriter) error {
	s.session.Options.MaxAge = -1
	return s.session.Save(r, w)
}

func (s *defaultSession) IsAuthenticated() bool {
	_, ok := s.session.Values["user_id"]
	return ok
}
