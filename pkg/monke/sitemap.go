package monke

import (
	"encoding/xml"
	"fmt"
	"html"
)

type urlEntry struct {
	Loc      string  `xml:"loc"`
	Lastmod  string  `xml:"lastmod"`
	Priority float32 `xml:"priority"`
}

type urlsetEntry struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	Urls    []urlEntry `xml:"url"`
}

func SitemapXml(urlPrefix string) ([]byte, error) {
	articles := Db.Articles.GetArticles(ArticleFilter{"", "", 0}, 0, 0)
	urls := make([]urlEntry, 0, len(articles))
	for _, article := range articles {
		articleLoc := html.EscapeString(SanitizeUrl(urlPrefix + article.Url))
		urls = append(urls, urlEntry{
			Loc:      articleLoc,
			Lastmod:  article.Created.Format("2006-01-02"),
			Priority: 0.8,
		})

		extras := article.Extras
		for k := range extras {
			urls = append(urls, urlEntry{
				Loc:      html.EscapeString(SanitizeUrl(fmt.Sprintf("%s/extra/%s/", articleLoc, k))),
				Lastmod:  article.Created.Format("2006-01-02"),
				Priority: 0.5,
			})
		}
	}
	urlset := urlsetEntry{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	return xml.MarshalIndent(urlset, "", "  ")

}
