package pages

import (
	"fmt"
	"html/template"
	"os"

	"github.com/labstack/echo/v4"

	"mrdmitry/blog/pkg/monke"
)

type ArticlePageData struct {
	HeadSnippet

	Nav  monke.NavData
	Body template.HTML
}

func ArticlePage(c echo.Context) error {
	path := c.Request().URL.Path
	readme := fmt.Sprintf("./web/data/%s/%s/README.md", c.Param("category"), c.Param("article"))
	var body []byte = nil

	body, err := monke.RenderMarkdownToHTMLAbs(readme, path)

	if err != nil {
		return NotFoundPage(c)
	}

	return c.Render(200, "article.html", ArticlePageData{
		HeadSnippet: NewHeadSnippet("Article"),
		Nav:         monke.Nav,
		Body:        template.HTML(string(body)),
	})
}

func ArticleAsset(c echo.Context) error {
	path := fmt.Sprintf("./web/data/%s/%s/assets/%s", c.Param("category"), c.Param("article"), c.Param("asset"))
	file, err := os.Open(path)
	if err != nil {
		return c.NoContent(404)
	}
	defer file.Close()
	return c.File(path)
}
