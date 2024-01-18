package monke

import (
	"encoding/xml"
	"fmt"
	"html"
	"time"

	"github.com/labstack/gommon/log"
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

func getLastmodTime(path string, fallback time.Time) (time.Time, error) {
	timestamp, err := GitLastLogTimeISO(path)
	if err != nil {
		log.Errorf("Could not find last log time for %s: %+v; falling back to: %v", path, err, fallback)
		return fallback, err
	}
	return timestamp, nil
}

func SitemapXml(urlPrefix string) ([]byte, error) {
	articles := Db.Articles.GetArticles(ArticleFilter{"", "", 0}, 0, 0)
	urls := make([]urlEntry, 0, len(articles)+2)

	// static entries
	fallbackTimestamp, err := getLastmodTime(".", time.Time{})
	if err != nil {
		panic("could not calculate fallback timestamp")
	}

	// home page
	timestamp, _ := getLastmodTime("web/templates/index.html", fallbackTimestamp)
	urls = append(urls, urlEntry{
		Loc:      html.EscapeString(SanitizeUrl(urlPrefix)),
		Lastmod:  timestamp.Format("2006-01-02"),
		Priority: 0.8,
	})

	// about page
	timestamp, _ = getLastmodTime("web/templates/about.html", fallbackTimestamp)
	urls = append(urls, urlEntry{
		Loc:      html.EscapeString(SanitizeUrl(urlPrefix + "/about/")),
		Lastmod:  timestamp.Format("2006-01-02"),
		Priority: 0.8,
	})

	// articles and their extras
	for _, article := range articles {
		timestamp, _ := getLastmodTime(article.ReadmePath, article.Created)
		articleLoc := html.EscapeString(SanitizeUrl(urlPrefix + article.Url))
		urls = append(urls, urlEntry{
			Loc:      articleLoc,
			Lastmod:  timestamp.Format("2006-01-02"),
			Priority: 0.8,
		})

		extras := article.Extras
		for k := range extras {
			timestamp, _ := getLastmodTime(fmt.Sprintf("%s/extra/%s/README.md", article.Root, k), article.Created)
			urls = append(urls, urlEntry{
				Loc:      html.EscapeString(SanitizeUrl(fmt.Sprintf("%s/extra/%s/", articleLoc, k))),
				Lastmod:  timestamp.Format("2006-01-02"),
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
