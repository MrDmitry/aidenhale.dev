package main

import (
    "html/template"
    "io"
    "log"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    monke "mrdmitry/blog/pkg/monke"
    pages "mrdmitry/blog/pkg/pages"
)

type BlogRenderer struct {
    templates *template.Template
}

func (t *BlogRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
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
    e.Renderer = &BlogRenderer {
        templates: tmpls,
    }

    e.Use(middleware.Logger())

    e.File("/favicon.ico", "./web/assets/favicon.ico")

    e.Static("/css", "./web/css")
    e.Static("/js", "./web/js")

    e.GET("/", pages.Index)

    e.Logger.Fatal(e.Start(":31337"))
}
