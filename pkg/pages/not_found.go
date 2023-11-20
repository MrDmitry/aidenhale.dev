package pages

import (
	"github.com/labstack/echo/v4"
)

func NotFoundPage(c echo.Context) error {
	return StaticPage(404, "404.html")(c)
}
