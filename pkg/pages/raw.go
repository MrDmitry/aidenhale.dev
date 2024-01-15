package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type RawProviderFunc func() ([]byte, error)

func RawXML(c echo.Context, provider RawProviderFunc) error {
	data, err := provider()
	if err != nil {
		log.Warnf("failed to generate XML response for %+v: %+v", c.Request().URL, err)
		return StaticPage(500, "500.html")(c)
	}
	return c.XMLBlob(200, data)
}
