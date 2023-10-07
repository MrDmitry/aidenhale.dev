package main

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"
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
		if parts[1] == k[1:] {
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

func ResolveRelativePaths() echo.MiddlewareFunc {
	// taken from https://github.com/labstack/echo/blob/4bc3e475e3137b6402933eec5e6fde641e0d2320/middleware/slash.go#L123-L130
	var sanitizeUri = func(uri string) string {
		// double slash `\\`, `//` or even `\/` is absolute uri for browsers and by redirecting request to that uri
		// we are vulnerable to open redirect attack. so replace all slashes from the beginning with single slash
		if len(uri) > 1 && (uri[0] == '\\' || uri[0] == '/') && (uri[1] == '\\' || uri[1] == '/') {
			uri = "/" + strings.TrimLeft(uri, `/\`)
		}
		return uri
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			url := req.URL
			path := url.Path
			if len(path) > 0 {
				absPath, err := filepath.Abs(path)
				if err != nil {
					log.Warnf("failed to process absolute path for %s: %+v", path, err)
					return c.NoContent(http.StatusBadRequest)
				}
				if path[len(path)-1] == '/' {
					absPath += "/"
				}
				absPath = sanitizeUri(absPath)
				log.Warnf("%s == %s", absPath, path)
				if absPath != path {
					return c.Redirect(http.StatusMovedPermanently, absPath)
				}
				// Forward
				req.RequestURI = absPath + "?" + c.QueryString()
				url.Path = absPath
			}
			return next(c)
		}
	}
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
	e.Pre(ResolveRelativePaths())
	e.Pre(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
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
