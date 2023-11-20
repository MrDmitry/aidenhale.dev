package pages

import "github.com/labstack/echo/v4"

func StaticPage(code int, entry string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(code, entry, nil)
	}
}
