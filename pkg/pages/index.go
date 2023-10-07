package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type IndexData struct {
	PageTitle string
	Nav       monke.NavData
}

func Index(c echo.Context) error {
	return c.Render(200, "index.html", IndexData{
		PageTitle: "Blog",
		Nav:       monke.Nav,
	})
}
