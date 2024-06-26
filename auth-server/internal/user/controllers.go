package user

import (
	"auth-server/internal"
	"auth-server/internal/common"
	"auth-server/internal/db"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewControllers(routes map[string]string, queries *db.Queries) *Controllers {
	return &Controllers{
		Routes: routes,
		userRepo: &defaultUserRepository{
			queries: queries,
		},
	}
}

type Controllers struct {
	Routes   map[string]string
	userRepo userRepository
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
	// TODO: check there's a current session
	// and don't attempt to login. Validate that the current
	// session is not expired.

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

	errs := form.Validate()
	if len(errs) > 0 {
		data.Errors = errs
		data.Email = form.Email
		return c.Render(http.StatusBadRequest, "login", data)
	}

	user, err := co.userRepo.Login(form.Email, form.Password)
	if err != nil {
		c.Logger().Error(err)
		if errors.Is(err, loginError) {
			data.Errors = map[string]string{
				"form": "Invalid email or password",
			}

			return c.Render(http.StatusBadRequest, "login", data)
		}
		data.Errors = map[string]string{
			"form": "Server Error",
		}
		return c.Render(http.StatusBadRequest, "login", data)
	}

	s, _ := common.NewDefaultSession(c.Request())
	s.SetUserInfo(common.UserInfo{
		UserID: user.ID,
		Email:  user.Email,
	})

	if err := s.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("Unable to save session: %w", err)
	}

	redirect := c.QueryParam("redirect")
	var rUrl string
	if len(redirect) > 0 {
		rUrl = redirect
	} else {
		rUrl = "/"
	}

	return c.Redirect(http.StatusMovedPermanently, rUrl)
}

func (co *Controllers) Logout(c echo.Context) error {
	s, _ := c.(common.AppContext).GetSession()

	if err := s.Delete(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("Unable to delete session: %w", err)
	}

	return c.Redirect(http.StatusMovedPermanently, co.Routes["index"])
}

type SignupData struct {
	internal.TemplateData
	Errors map[string]string
	Email  string
}

func (co *Controllers) Signup(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", SignupData{
		TemplateData: internal.TemplateData{Routes: co.Routes},
	})
}

func (co *Controllers) HandleSignupForm(c echo.Context) error {
	data := SignupData{
		TemplateData: internal.TemplateData{Routes: co.Routes},
	}

	var form signupForm
	if err := c.Bind(&form); err != nil {
		data.Errors = map[string]string{
			"form": "Invalid form data",
		}

		return c.Render(http.StatusInternalServerError, "signup", data)
	}

	errs := form.Validate()
	if len(errs) > 0 {
		data.Errors = errs
		data.Email = form.Email
		return c.Render(http.StatusBadRequest, "signup", data)
	}

	err := co.userRepo.Create(form.Email, form.Password)
	if err != nil {
		code := http.StatusInternalServerError
		msg := "Server Error. Try again later."
		if errors.Is(err, userEmailTaken) {
			code = http.StatusBadRequest
			msg = "Email is already taken"
		}

		data.Errors = map[string]string{
			"form": msg,
		}

		c.Logger().Error(err)

		return c.Render(code, "signup", data)
	}

	return c.Redirect(http.StatusMovedPermanently, co.Routes["login"])
}
