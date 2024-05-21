package internal

import (
	"fmt"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

var layoutPath = "views/layout.html"
var templates = map[string]string{
	"index":    "views/index.html",
	"register": "views/register.html",
}

type Templates struct {
	templates map[string]*template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("Template %s not found", name)
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

func ParseTemplates() *Templates {
	base := template.Must(template.ParseFiles(layoutPath))
	t := make(map[string]*template.Template, len(templates))

	for name, path := range templates {
		t[name] = template.Must(template.Must(base.Clone()).ParseFiles(path))
	}

	return &Templates{
		templates: t,
	}
}

type TemplateData struct {
	Routes map[string]string
}
