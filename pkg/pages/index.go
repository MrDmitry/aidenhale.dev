package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type IndexPageData struct {
	HeadSnippet
	ArticlesSnippetData

	Nav monke.NavData
}

func IndexPage(c echo.Context) error {
	return c.Render(200, "index.html", IndexPageData{
		HeadSnippet:         NewHeadSnippet("Blog"),
		ArticlesSnippetData: NewArticlesSnippetData(c, monke.ArticleFilter{}),
		Nav:                 monke.Nav,
	})
}
