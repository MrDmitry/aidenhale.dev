package pages

import (
	"mrdmitry/blog/pkg/monke"

	"github.com/labstack/echo/v4"
)

type TagTally struct {
	Name  string
	Count int
}

type TagsPageData struct {
	HeadSnippet

	Nav  monke.NavData
	Tags []TagTally
}

func TagsPage(c echo.Context) error {
	tags := monke.Db.Articles.GetTagsSizes()
	tagTally := make([]TagTally, 0, len(tags))
	for tag, count := range tags {
		tagTally = append(tagTally, TagTally{
			Name:  tag,
			Count: count,
		})
	}
	return c.Render(200, "tags.html", TagsPageData{
		HeadSnippet: NewHeadSnippet("Tags"),
		Nav:         monke.Nav,
		Tags:        tagTally,
	})
}
