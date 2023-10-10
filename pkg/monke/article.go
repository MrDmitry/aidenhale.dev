package monke

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

// struct passed to the html template
type Article struct {
	Root       string // filepath to the article root location
	Id         string // human readable identifier
	Url        string // article path
	ReadmePath string // filepath to README.md with article contents
	Summary    string // plain text summary

	ArticleMeta
}

// struct for toml unmarhsalling
type ArticleMeta struct {
	Title   string
	Created time.Time
}

func NewArticle(f string, c *Category) (*Article, error) {
	dir, err := os.Open(f)
	if err != nil {
		log.Fatalf("Failed to open %s: %+v", f, err)
		return nil, err
	}
	defer dir.Close()

	_, err = dir.ReadDir(0)
	if err != nil {
		log.Fatalf("Failed to read %s: %+v", f, err)
		return nil, err
	}

	var meta ArticleMeta

	monkeMeta := path.Join(f, "monke.toml")
	tree, err := toml.LoadFile(monkeMeta)

	if err != nil {
		log.Warnf("could not process %s, skipping: %+v", monkeMeta, err)
		return nil, err
	}

	err = tree.Unmarshal(&meta)

	if err != nil {
		log.Warnf("could not unwrap %s, skipping: %+v", monkeMeta, err)
		return nil, err
	}

	readme := path.Join(f, "README.md")
	summary, err := RenderMarkdownToTextPreview(readme, 200)

	if err != nil {
		log.Warnf("could not process %s, skipping: %+v", readme, err)
		return nil, err
	}

	article := new(Article)
	article.Root = dir.Name()
	article.Id = filepath.Base(article.Root)
	article.Url = sanitizeUrl(fmt.Sprintf("%s/%s/", c.Url, article.Id))
	article.ReadmePath = readme
	article.Summary = string(summary)
	article.ArticleMeta = meta

	return article, nil
}

func sanitizeUrl(s string) string {
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	return s + "/"
}
