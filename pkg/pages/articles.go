package pages

import (
	"fmt"
	"mrdmitry/blog/pkg/monke"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type annotatedArticle struct {
	Article *monke.Article
	IsLast  bool
	NextUrl string
	Tags    TagData
}

type ArticlesSnippetData struct {
	Articles []annotatedArticle
	Tags     TagData
	Filter   monke.ArticleFilter
}

func NewArticlesSnippetData(c echo.Context, f monke.ArticleFilter) ArticlesSnippetData {
	page := f.Page
	limit := 4

	articlesSrc := monke.Db.Articles.GetArticles(f, limit, limit*page)
	if len(articlesSrc) == 0 {
		return ArticlesSnippetData{Articles: nil, Filter: f}
	}

	articles := make([]annotatedArticle, 0, len(articlesSrc))
	for _, a := range articlesSrc {
		articles = append(articles, annotatedArticle{
			Article: a,
			Tags:    NewArticleTagData(a, f.Tag, false),
		})
	}

	lastArticle := &articles[len(articles)-1]
	lastArticle.IsLast = true

	nextUrl := f.ToUrlValues()
	nextUrl.Set("page", strconv.Itoa(page+1))

	lastArticle.NextUrl = fmt.Sprintf("/articles/?%s", nextUrl.Encode())

	return ArticlesSnippetData{
		Articles: articles,
		Tags:     NewGlobalTagData(f.Tag, false),
		Filter:   f,
	}
}

func ArticlesSnippet(c echo.Context) error {
	var filter monke.ArticleFilter
	err := c.Bind(&filter)
	if err != nil {
		log.Warnf("Failed to bind filter parameters: %+v", c.Request().URL)
		return c.NoContent(400)
	}
	return c.Render(200, "articles.html", NewArticlesSnippetData(c, filter))
}
