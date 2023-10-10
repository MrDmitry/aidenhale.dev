package pages

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"mrdmitry/blog/pkg/monke"
)

type BlogData struct {
	PageTitle string
	Nav       monke.NavData
	Body      template.HTML
	Articles  []*monke.Article
}

func Blog(c echo.Context) error {
	categoryId := c.Param("category")

	if monke.Db.Categories[categoryId] == nil {
		return NotFound(c)
	}

	category := monke.Db.Categories[categoryId]
	body, _ := monke.RenderMarkdownToHTML(category.ReadmePath)
	articles := category.GetArticlesByTime(10, 0)

	return c.Render(200, "blog.html", BlogData{
		PageTitle: "Articles",
		Nav:       monke.Nav,
		Body:      template.HTML(string(body)),
		Articles:  articles,
	})
}
