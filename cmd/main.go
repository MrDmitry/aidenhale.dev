package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	monke "mrdmitry/blog/pkg/monke"
	pages "mrdmitry/blog/pkg/pages"
)

type Background string

const (
	Dark  Background = "dark"
	Light            = "light"
)

func (bg Background) String() string {
	return string(bg)
}

type TemplateEntry struct {
	tmpl *template.Template
	main string
	bg   Background
}

func newUrl(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Errorf("failed to create url from %+v", s)
		return nil
	}
	return u
}

func atoi(s string) int {
	v, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return v
}

func generateUrl(u *url.URL, t string) string {
	values := u.Query()
	if len(t) > 0 {
		values.Set("tag", t)
	} else {
		values.Del("tag")
	}
	values.Del("page")
	if len(values) > 0 {
		return fmt.Sprintf("%s?%s", u.Path, values.Encode())
	} else {
		return u.Path
	}
}

func newTemplateEntry(ts []string, m string, bg Background) TemplateEntry {
	return TemplateEntry{
		tmpl: template.Must(template.New(m).Funcs(template.FuncMap{
			"newUrl":      newUrl,
			"int":         atoi,
			"generateUrl": generateUrl,
		}).ParseFiles(ts...)),
		main: m,
		bg:   bg,
	}

}

type BlogRenderer struct {
	templates map[string]TemplateEntry
}

type dataWrapper struct {
	pages.HeadSnippet

	Background string
	Url        *url.URL
	Data       interface{}
}

func (r *BlogRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	entry, ok := r.templates[name]
	if !ok {
		log.Errorf("failed to find template %s", name)
		return c.NoContent(500)
	}
	return entry.tmpl.ExecuteTemplate(w, entry.main, dataWrapper{
		HeadSnippet: pages.NewHeadSnippet(),
		Background:  entry.bg.String(),
		Url:         c.Request().URL,
		Data:        data,
	})
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	c.Render(code, errorPage, nil)
}

var staticDirs = map[string]string{
	"/assets": "./web/assets",
	"/css":    "./web/dist/css",
	"/js":     "./web/js",
}

var staticFiles = map[string]string{
	"/favicon.ico": "./web/assets/favicon.ico",
}

var generatedFiles = map[string]string{
	"/robots.txt":  "",
	"/sitemap.xml": "",
}

func AssetSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return true
	}

	collections := []map[string]string{staticFiles, generatedFiles, staticDirs}
	for _, col := range collections {
		for k := range col {
			if parts[1] == k[1:] {
				return true
			}
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
	protocol := "https"
	flag.Func("protocol", "webserver protocol: [http, https] (default \"https\")", func(value string) error {
		switch value {
		case "http":
			fallthrough
		case "https":
			protocol = value
			return nil
		default:
			return errors.New("unexpected value, expected one of [http, https]")
		}
	})
	hostname := flag.String("hostname", "aidenhale.dev", "hostname")
	portPtr := flag.Int("port", 31337, "port to bind to")
	flag.Parse()

	urlPrefix := monke.SanitizeUrl(fmt.Sprintf("%s://%s/", protocol, *hostname))

	tmpls := make(map[string]TemplateEntry)

	snippets := []string{
		"./web/templates/articleCard.html",
		"./web/templates/tagLabels.html",
	}

	tmpls["404.html"] = newTemplateEntry(
		[]string{
			"./web/templates/404.html",
			"./web/templates/base.html",
		},
		"base.html",
		Light,
	)

	tmpls["500.html"] = newTemplateEntry(
		[]string{
			"./web/templates/500.html",
			"./web/templates/base.html",
		},
		"base.html",
		Light,
	)

	tmpls["article.html"] = newTemplateEntry(
		append([]string{
			"./web/templates/article.html",
			"./web/templates/base.html",
		}, snippets...),
		"base.html",
		Dark,
	)

	tmpls["articleExtra.html"] = newTemplateEntry(
		append([]string{
			"./web/templates/articleExtra.html",
			"./web/templates/base.html",
		}, snippets...),
		"base.html",
		Light,
	)

	tmpls["articles.html"] = newTemplateEntry(
		append([]string{
			"./web/templates/articles.html",
			"./web/templates/articlesIndex.html",
		}, snippets...),
		"articlesIndex.html",
		Light,
	)

	tmpls["index.html"] = newTemplateEntry(
		append([]string{
			"./web/templates/index.html",
			"./web/templates/articles.html",
			"./web/templates/base.html",
		}, snippets...),
		"base.html",
		Light,
	)

	tmpls["about.html"] = newTemplateEntry(
		append([]string{
			"./web/templates/about.html",
			"./web/templates/base.html",
		}, snippets...),
		"base.html",
		Light,
	)

	tmpls["robots.txt"] = newTemplateEntry(
		[]string{
			"./web/templates/robots.txt",
		},
		"robots.txt",
		"",
	)

	err := monke.InitDb("./web/data")

	if err != nil {
		log.Fatalf("could not initialize database: %+v", err)
	}

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
	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper:      func(c echo.Context) bool { return !AssetSkipper(c) },
		RedirectCode: http.StatusMovedPermanently,
	}))

	for k, v := range staticFiles {
		e.File(k, v)
	}

	for k, v := range staticDirs {
		e.Static(k, v)
	}

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.GET("/", pages.IndexPage)
	e.GET("/sitemap.xml", func(c echo.Context) error {
		return pages.RawXML(c, func() ([]byte, error) {
			return monke.SitemapXml(urlPrefix)
		})
	})
	e.GET("/robots.txt", func(c echo.Context) error {
		return c.Render(200, "robots.txt", struct {
			UrlPrefix string
		}{
			UrlPrefix: urlPrefix,
		})
	})
	e.GET("/about/", pages.StaticPage(200, "about.html"))
	e.GET("/articles/", pages.ArticlesSnippet)
	e.GET("/blog/:category/:article/", pages.ArticlePage)
	e.GET("/blog/:category/:article/extra/:extra/", pages.ArticleAppendixPage)
	e.GET("/blog/:category/:article/assets/:asset", pages.ArticleAsset)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(*portPtr)))
}
