package pages

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"mrdmitry/blog/pkg/monke"
)

type BlogData struct {
	PageTitle string
	Nav       monke.NavData
	Body      template.HTML
	Articles  []monke.Article
}

func Blog(c echo.Context) error {
	topic := c.Param("topic")
	readme := fmt.Sprintf("./web/data/%s/README.md", topic)
	var body []byte = nil

	body, err := monke.RenderMarkdownToHTML(readme)

	if err != nil {
		log.Warnf("could not generate body: %+v", readme, err)
		body = nil
	}

	articles, _ := monke.GetArticles(topic)

	return c.Render(200, "blog.html", BlogData{
		PageTitle: "Articles",
		Nav:       monke.Nav,
		Body:      template.HTML(string(body)),
		Articles:  articles,
	})
}
