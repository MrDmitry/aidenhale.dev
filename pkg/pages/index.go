package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type IndexPageData struct {
	ArticlesSnippetData
}

func IndexPage(c echo.Context) error {
	var filter monke.ArticleFilter
	err := c.Bind(&filter)
	if err != nil {
		log.Warnf("Failed to bind filter parameters: %+v", c.Request().URL)
		filter = monke.ArticleFilter{}
	}
	return c.Render(200, "index.html", IndexPageData{
		ArticlesSnippetData: NewArticlesSnippetData(c, filter),
	})
}
