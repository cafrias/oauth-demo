package main

import (
	"auth-server/internal"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	Routes map[string]string
}

func (h *Handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", internal.TemplateData{Routes: h.Routes})
}

type AppInfo struct {
	Name         string
	Type         string
	RedirectURI  string
	ClientId     string
	ClientSecret string
}

type RegisterData struct {
	internal.TemplateData
	Success bool
	Errors  map[string]string
	AppInfo AppInfo
}

func (h *Handler) Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", RegisterData{TemplateData: internal.TemplateData{Routes: h.Routes}})
}

type RegisterForm struct {
	Name        string `form:"name"`
	Type        string `form:"type"`
	RedirectURI string `form:"redirect_uri"`
}

func (f *RegisterForm) Validate() map[string]string {
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

func (h *Handler) HandleRegisterForm(c echo.Context) error {
	data := RegisterData{
		TemplateData: internal.TemplateData{Routes: h.Routes},
	}
	fmt.Println("OK!")

	var form RegisterForm
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
		data.AppInfo = AppInfo{
			Name:        form.Name,
			Type:        form.Type,
			RedirectURI: form.RedirectURI,
		}
		return c.Render(http.StatusBadRequest, "register", data)
	}

	// TODO: include logic to generate client ID and client secret
	clientId := "123456"
	var clientSecret string
	if form.Type == "server-side" {
		clientSecret = "123456"
	}

	data.Success = true
	data.AppInfo = AppInfo{
		Name:         form.Name,
		Type:         form.Type,
		RedirectURI:  form.RedirectURI,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	return c.Render(http.StatusOK, "register", data)
}

func CreateRoutes() map[string]string {
	return map[string]string{
		"index":        "/",
		"app/register": "/register",
	}
}

func main() {
	templates := internal.ParseTemplates()
	e := echo.New()
	e.Renderer = templates
	e.Use(middleware.Logger())
	routes := CreateRoutes()
	h := &Handler{
		Routes: routes,
	}

	e.GET(routes["index"], h.Index)
	e.GET(routes["app/register"], h.Register)
	e.POST(routes["app/register"], h.HandleRegisterForm)

	e.Logger.Fatal(e.Start(":1323"))
}
