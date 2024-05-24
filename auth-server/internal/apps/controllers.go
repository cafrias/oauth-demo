package apps

import (
	"auth-server/internal"
	"auth-server/internal/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewControllers(routes map[string]string) *Controllers {
	return &Controllers{
		Routes:        routes,
		appRepository: &defaultAppRepository{},
	}
}

type Controllers struct {
	Routes        map[string]string
	appRepository appRepository
}

type RegisterData struct {
	internal.TemplateData
	Success bool
	Errors  map[string]string
	AppInfo App
}

func (co *Controllers) Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", RegisterData{TemplateData: internal.TemplateData{Routes: co.Routes}})
}

func (co *Controllers) HandleRegisterForm(c echo.Context) error {
	data := RegisterData{
		TemplateData: internal.TemplateData{Routes: co.Routes},
	}

	var form registerForm
	if err := c.Bind(&form); err != nil {
		data.Errors = map[string]string{
			"form": "Invalid form data",
		}

		return c.Render(http.StatusBadRequest, "register", data)
	}

	// TODO: add validation for the form fields
	errors := form.Validate()
	if len(errors) > 0 {
		data.Success = false
		data.Errors = errors
		data.AppInfo = App{
			Name:        form.Name,
			Type:        form.Type,
			RedirectURI: form.RedirectURI,
		}
		return c.Render(http.StatusBadRequest, "register", data)
	}

	// TODO: include logic to generate client ID and client secret
	s, _ := c.(common.AppContext).GetSession()
	u := s.GetUserInfo()

	app, err := co.appRepository.Register(registerInput{
		UserID:      u.UserID,
		Name:        form.Name,
		Type:        form.Type,
		RedirectURI: form.RedirectURI,
	})
	if err != nil {
		// TODO: handle error
		return err
	}

	data.Success = true
	data.AppInfo = *app

	return c.Render(http.StatusOK, "register", data)
}
