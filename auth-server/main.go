package main

import (
	"auth-server/internal"
	"auth-server/internal/apps"
	"auth-server/internal/auth"
	"auth-server/internal/common"
	"auth-server/internal/db"
	"auth-server/internal/user"
	"database/sql"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
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
		"apps/list":     "/apps/list",
		"apps/delete":   "/apps/delete",
		"apps/reset":    "/apps/reset",
		"login":         "/login",
		"logout":        "/logout",
		"signup":        "/signup",
	}
}

func main() {
	conn, err := sql.Open("sqlite3", os.Getenv("AUTH_SERVER_DB_PATH"))
	if err != nil {
		panic(err)
	}

	queries := db.New(conn)

	templates := internal.ParseTemplates()
	e := echo.New()
	e.Use(common.UseAppContext)
	e.Use(middleware.Logger())
	e.Renderer = templates
	routes := CreateRoutes()
	h := &Handler{
		Routes: routes,
	}

	a := apps.NewControllers(routes, queries)
	appRoutes := e.Group("/apps")
	appRoutes.Use(auth.Authenticated)
	appRoutes.GET("/register", a.Register)
	appRoutes.POST("/register", a.HandleRegisterForm)
	appRoutes.POST("/delete", a.DeleteApp)
	appRoutes.POST("/reset", a.ResetAppSecret)
	appRoutes.GET("/list", a.UserApps)

	u := user.NewControllers(routes, queries)
	e.GET(routes["login"], u.Login)
	e.POST(routes["login"], u.HandleLoginForm)
	e.POST(routes["logout"], u.Logout)
	e.GET(routes["signup"], u.Signup)
	e.POST(routes["signup"], u.HandleSignupForm)

	e.GET(routes["index"], h.Index)

	e.Logger.Fatal(e.Start(":1323"))
}
