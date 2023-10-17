package monke

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

// struct passed to the html template
type Article struct {
	Category   string // article category identifier
	Root       string // filepath to the article root location
	Id         string // human readable identifier
	Url        string // article path
	ReadmePath string // filepath to README.md with article contents
	Summary    string // plain text summary
	WordCount  int

	ArticleData
}

type articleMeta struct {
	Version int
}

type ArticleData struct {
	Title   string
	Created time.Time
	Tags    []string
}

// struct for toml unmarshalling
type ArticleToml struct {
	Meta articleMeta
	Data ArticleData
}

func NewArticle(f string, c string, urlPrefix string) (*Article, error) {
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

	var meta ArticleToml

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

	if meta.Meta.Version != 0 {
		log.Warnf("incompatible version detected: expected %i, found %i", 0, meta.Meta.Version)
		return nil, err
	}

	readme := path.Join(f, "README.md")
	summary, err := RenderMarkdownToText(readme)

	if err != nil {
		log.Warnf("could not process %s, skipping: %+v", readme, err)
		return nil, err
	}

	wordCount := len(strings.Fields(string(summary)))

	threshold := min(len(summary), 200)
	if threshold != len(summary) {
		summary = summary[:threshold-3]
		summary = append(summary, []byte("...")...)
	}

	article := new(Article)
	article.Category = c
	article.Root = dir.Name()
	article.Id = filepath.Base(article.Root)
	article.Url = sanitizeUrl(fmt.Sprintf("%s/%s/", urlPrefix, article.Id))
	article.ReadmePath = readme
	article.Summary = string(summary)
	article.ArticleData = meta.Data
	article.WordCount = wordCount

	return article, nil
}

func sanitizeUrl(s string) string {
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	return s + "/"
}
