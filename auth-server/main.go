package main

import (
	"auth-server/internal"
	"auth-server/internal/apps"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	Routes map[string]string
}

func (h *Handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", internal.TemplateData{Routes: h.Routes})
}

func CreateRoutes() map[string]string {
	return map[string]string{
		"index":         "/",
		"apps/register": "/apps/register",
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

	appsController := apps.NewControllers(routes)

	e.GET(routes["index"], h.Index)
	e.GET(routes["apps/register"], appsController.Register)
	e.POST(routes["apps/register"], appsController.HandleRegisterForm)

	e.Logger.Fatal(e.Start(":1323"))
}
