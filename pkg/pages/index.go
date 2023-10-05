package pages

import (
    "github.com/labstack/echo/v4"

    "mrdmitry/blog/pkg/monke"
)

type IndexData struct {
    Nav monke.NavData
}

func Index(c echo.Context) error {
    return c.Render(200, "index.html", IndexData{
        Nav: monke.Nav,
    })
}
