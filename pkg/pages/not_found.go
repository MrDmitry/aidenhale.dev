package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type NotFoundPageData struct {
	HeadSnippet

	Nav monke.NavData
}

func NotFoundPage(c echo.Context) error {
	return c.Render(404, "404.html", NotFoundPageData{
		HeadSnippet: NewHeadSnippet("Not found"),
		Nav:         monke.Nav,
	})
}
