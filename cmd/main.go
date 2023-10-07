package main

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	monke "mrdmitry/blog/pkg/monke"
	pages "mrdmitry/blog/pkg/pages"
)

type BlogRenderer struct {
	templates *template.Template
}

func (t *BlogRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var statics = map[string]string{
	"/assets": "./web/assets",
	"/css":    "./web/css",
	"/js":     "./web/js",
}

func AssetSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	parts := strings.Split(path, "/")
	for k, _ := range statics {
		if parts[0] == k[1:] {
			return true
		}
	}
	for _, part := range parts {
		if part == "assets" {
			return true
		}
	}
	return false
}

func main() {
	tmpls, err := template.ParseGlob(
		"./web/templates/*.html",
	)

	if err != nil {
		log.Fatalf("could not load templates: %+v", err)
	}

	monke.NavInit()

	e := echo.New()
	e.Renderer = &BlogRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper:      AssetSkipper,
		RedirectCode: http.StatusMovedPermanently,
	}))

	e.File("/favicon.ico", "./web/assets/favicon.ico")

	for k, v := range statics {
		e.Static(k, v)
	}

	e.GET("/", pages.Index)
	e.GET("/blog/:topic/", pages.Blog)
	e.GET("/blog/:topic/:article/", pages.Article)
	e.GET("/blog/:topic/:article/assets/:asset", pages.ArticleAssets)

	e.Logger.Fatal(e.Start(":31337"))
}
