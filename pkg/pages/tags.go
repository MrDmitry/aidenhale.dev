package pages

import "mrdmitry/blog/pkg/monke"

type TagData struct {
	Tags       []string
	CurrentTag string
	Clickable  bool
}

func NewArticleTagData(article *monke.Article, t string, c bool) TagData {
	return TagData{
		Tags:       article.Tags,
		CurrentTag: t,
		Clickable:  c,
	}
}

func NewGlobalTagData(t string, c bool) TagData {
	return TagData{
		Tags:       monke.Db.Articles.GetTags(),
		CurrentTag: t,
		Clickable:  c,
	}
}
