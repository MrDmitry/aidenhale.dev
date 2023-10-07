package pages

import (
	"fmt"
	"html/template"
	"os"

	"github.com/labstack/echo/v4"

	"mrdmitry/blog/pkg/monke"
)

type ArticleData struct {
	PageTitle string
	Nav       monke.NavData
	Body      template.HTML
}

func Article(c echo.Context) error {
	path := c.Request().URL.Path
	readme := fmt.Sprintf("./web/data/%s/%s/README.md", c.Param("topic"), c.Param("article"))
	var body []byte = nil

	body, err := monke.RenderMarkdownAbs(readme, path)

	if err != nil {
		return NotFound(c)
	}

	return c.Render(200, "article.html", ArticleData{
		PageTitle: "Article",
		Nav:       monke.Nav,
		Body:      template.HTML(string(body)),
	})
}

func ArticleAssets(c echo.Context) error {
	path := fmt.Sprintf("./web/data/%s/%s/assets/%s", c.Param("topic"), c.Param("article"), c.Param("asset"))
	file, err := os.Open(path)
	if err != nil {
		return c.NoContent(404)
	}
	defer file.Close()
	return c.File(path)
}
