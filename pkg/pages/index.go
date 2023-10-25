package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type IndexPageData struct {
	ArticlesSnippetData
}

func IndexPage(c echo.Context) error {
	return c.Render(200, "index.html", IndexPageData{
		ArticlesSnippetData: NewArticlesSnippetData(c, monke.ArticleFilter{}),
	})
}
