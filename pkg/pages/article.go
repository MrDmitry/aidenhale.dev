package pages

import (
	"fmt"
	"html/template"
	"os"

	"github.com/labstack/echo/v4"

	"mrdmitry/blog/pkg/monke"
)

type ArticlePageData struct {
	Article  *monke.Article
	Title    string
	Body     template.HTML
	Tags     TagData
	ExtraKey string
}

func ArticleAppendixPage(c echo.Context) error {
	path := c.Request().URL.Path
	article := monke.Db.Articles.GetArticle(c.Param("category"), c.Param("article"))
	if article == nil {
		return NotFoundPage(c)
	}

	extraKey := c.Param("extra")
	extra, ok := article.Extras[extraKey]
	if ok == false {
		return NotFoundPage(c)
	}

	readme := fmt.Sprintf("./web/data/%s/%s/extra/%s/README.md",
		c.Param("category"),
		c.Param("article"),
		extraKey,
	)
	var body []byte = nil

	body, err := monke.RenderMarkdownToHTMLAbs(readme, path)

	if err != nil {
		return NotFoundPage(c)
	}

	return c.Render(200, "articleExtra.html", ArticlePageData{
		Article:  article,
		Title:    extra.Title,
		Body:     template.HTML(string(body)),
		ExtraKey: extraKey,
	})
}

func ArticlePage(c echo.Context) error {
	path := c.Request().URL.Path
	article := monke.Db.Articles.GetArticle(c.Param("category"), c.Param("article"))
	if article == nil {
		return NotFoundPage(c)
	}

	readme := fmt.Sprintf("./web/data/%s/%s/README.md", c.Param("category"), c.Param("article"))
	var body []byte = nil

	body, err := monke.RenderMarkdownToHTMLAbs(readme, path)

	if err != nil {
		return NotFoundPage(c)
	}

	return c.Render(200, "article.html", ArticlePageData{
		Article: article,
		Title:   article.Title,
		Body:    template.HTML(string(body)),
		Tags:    NewArticleTagData(article, "", true),
	})
}

func ArticleAsset(c echo.Context) error {
	path := fmt.Sprintf("./web/data/%s/%s/assets/%s", c.Param("category"), c.Param("article"), c.Param("asset"))
	file, err := os.Open(path)
	if err != nil {
		return NotFoundPage(c)
	}
	defer file.Close()
	return c.File(path)
}
