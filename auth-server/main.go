package main

import (
	"auth-server/internal"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Routes map[string]string
}

func (h *Handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", internal.TemplateData{Routes: h.Routes})
}

func CreateRoutes() map[string]string {
	return map[string]string{
		"index":        "/",
		"app/register": "/register",
	}
}

func main() {
	templates := internal.NewTemplate("views/*.html")
	e := echo.New()
	e.Renderer = templates
	routes := CreateRoutes()
	h := &Handler{
		Routes: routes,
	}

	e.GET(routes["index"], h.Index)

	e.Logger.Fatal(e.Start(":1323"))
}
