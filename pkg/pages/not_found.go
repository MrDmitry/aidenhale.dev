package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type NotFoundData struct {
	PageTitle string
	Nav       monke.NavData
}

func NotFound(c echo.Context) error {
	return c.Render(404, "404.html", NotFoundData{
		PageTitle: "Not found",
		Nav:       monke.Nav,
	})
}
