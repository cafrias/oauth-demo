package user

import (
	"auth-server/internal"
	"auth-server/internal/common"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewControllers(routes map[string]string) *Controllers {
	return &Controllers{
		Routes: routes,
	}
}

type Controllers struct {
	Routes map[string]string
}

type LoginData struct {
	internal.TemplateData
	Errors map[string]string
	Email  string
}

func (co *Controllers) Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", LoginData{TemplateData: internal.TemplateData{Routes: co.Routes}})
}

func (co *Controllers) HandleLoginForm(c echo.Context) error {
	ctx := c.(common.AppContext)
	s, err := ctx.GetSession()
	if err != nil {
		return fmt.Errorf("Unable to parse session: %w", err)
	}

	data := LoginData{
		TemplateData: internal.TemplateData{Routes: co.Routes},
	}

	var form loginForm
	if err := c.Bind(&form); err != nil {
		data.Errors = map[string]string{
			"form": "Invalid form data",
		}

		return c.Render(http.StatusInternalServerError, "login", data)
	}

	errors := form.Validate()
	if len(errors) > 0 {
		data.Errors = errors
		data.Email = form.Email
		return c.Render(http.StatusBadRequest, "login", data)
	}

	s.SetUserInfo(common.UserInfo{
		UserID: "123",
		Email:  form.Email,
	})

	if err := ctx.SaveSession(s); err != nil {
		return fmt.Errorf("Unable to save session: %w", err)
	}

	return c.Redirect(http.StatusMovedPermanently, co.Routes["index"])
}

func (co *Controllers) Logout(c echo.Context) error {
	ctx := c.(common.AppContext)
	s, err := ctx.GetSession()
	if err != nil {
		return fmt.Errorf("Unable to parse session: %w", err)
	}

	if err := ctx.DeleteSession(s); err != nil {
		return fmt.Errorf("Unable to delete session: %w", err)
	}

	return c.Redirect(http.StatusMovedPermanently, co.Routes["index"])
}
