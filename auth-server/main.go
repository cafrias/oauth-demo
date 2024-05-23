package main

import (
	"auth-server/internal"
	"auth-server/internal/apps"
	"auth-server/internal/common"
	"auth-server/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	Routes map[string]string
}

type HomeData struct {
	internal.TemplateData
	Email string
}

func (h *Handler) Index(c echo.Context) error {
	ctx := c.(common.AppContext)

	s, err := ctx.GetSession()
	if err != nil {
		return err
	}

	data := HomeData{
		TemplateData: internal.TemplateData{Routes: h.Routes},
		Email:        s.GetUserInfo().Email,
	}

	return c.Render(http.StatusOK, "index", data)
}

func CreateRoutes() map[string]string {
	return map[string]string{
		"index":         "/",
		"apps/register": "/apps/register",
		"login":         "/login",
		"logout":        "/logout",
		"signup":        "/signup",
	}
}

func main() {
	templates := internal.ParseTemplates()
	e := echo.New()
	e.Use(common.UseAppContext)
	e.Use(middleware.Logger())
	e.Renderer = templates
	routes := CreateRoutes()
	h := &Handler{
		Routes: routes,
	}

	e.GET(routes["index"], h.Index)

	a := apps.NewControllers(routes)
	e.GET(routes["apps/register"], a.Register)
	e.POST(routes["apps/register"], a.HandleRegisterForm)

	u := user.NewControllers(routes)
	e.GET(routes["login"], u.Login)
	e.POST(routes["login"], u.HandleLoginForm)
	e.POST(routes["logout"], u.Logout)

	e.Logger.Fatal(e.Start(":1323"))
}
