package apps

import "net/url"

type registerForm struct {
	Name        string `form:"name"`
	Type        string `form:"type"`
	RedirectURI string `form:"redirect_uri"`
	CsrfToken   string `form:"csrf_token"`
}

func (f *registerForm) Validate() map[string]string {
	errors := make(map[string]string)

	if len(f.Name) == 0 || len(f.Name) >= 255 {
		errors["name"] = "Name is required and must be less than 255 characters"
	}

	switch f.Type {
	case "server-side":
	case "client-side":
	case "native":
	default:
		errors["type"] = "Invalid application type"
	}

	if len(f.RedirectURI) == 0 {
		errors["redirect_uri"] = "Redirect URI is required"
	} else if _, err := url.ParseRequestURI(f.RedirectURI); err != nil {
		errors["redirect_uri"] = "Invalid redirect URI"
	}

	return errors
}

type App struct {
	ID           string
	ClientID     string
	ClientSecret string
	UserID       string
	Name         string
	RedirectURI  string
	Type         string
}
