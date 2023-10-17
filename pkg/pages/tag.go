package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type TagPageData struct {
	HeadSnippet
	ArticlesSnippetData

	Nav monke.NavData
}

func TagPage(c echo.Context) error {
	tag := c.Param("tag")

	articles := NewArticlesSnippetData(c, monke.ArticleFilter{Tag: tag})

	if articles.Articles == nil {
		return NotFoundPage(c)
	}

	return c.Render(200, "tag.html", TagPageData{
		HeadSnippet:         NewHeadSnippet("Articles"),
		Nav:                 monke.Nav,
		ArticlesSnippetData: articles,
	})
}
