package pages

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"mrdmitry/blog/pkg/monke"
)

type CategoryPageData struct {
	ArticlesSnippetData

	Body template.HTML
}

func CategoryPage(c echo.Context) error {
	categoryId := c.Param("category")

	if monke.Db.Categories[categoryId] == nil {
		return NotFoundPage(c)
	}

	category := monke.Db.Categories[categoryId]
	body, _ := monke.RenderMarkdownToHTML(category.ReadmePath)

	return c.Render(200, "category.html", CategoryPageData{
		Body:                template.HTML(string(body)),
		ArticlesSnippetData: NewArticlesSnippetData(c, monke.ArticleFilter{Category: categoryId}),
	})
}
