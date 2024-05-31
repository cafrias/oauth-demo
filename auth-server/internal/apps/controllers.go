package apps

import (
	"auth-server/internal"
	"auth-server/internal/common"
	"auth-server/internal/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewControllers(routes map[string]string, queries *db.Queries) *Controllers {
	return &Controllers{
		Routes: routes,
		appRepository: &defaultAppRepository{
			queries: queries,
		},
	}
}

type Controllers struct {
	Routes        map[string]string
	appRepository appRepository
}

type RegisterData struct {
	internal.TemplateData
	Errors   map[string]string
	FormData registerForm
	Result   *App
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
		data.Errors = errors
		data.FormData = form
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

	data.Result = app

	return c.Render(http.StatusOK, "register", data)
}

type UserAppsData struct {
	internal.TemplateData
	Apps []App
}

func (co *Controllers) UserApps(c echo.Context) error {
	u, _ := c.(common.AppContext).GetUser()

	apps, err := co.appRepository.GetAllByUser(u.UserID)
	if err != nil {
		return err
	}

	data := UserAppsData{
		TemplateData: internal.TemplateData{Routes: co.Routes},
		Apps:         apps,
	}
	return c.Render(http.StatusOK, "user-apps", data)
}
