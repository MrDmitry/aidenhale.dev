package monke

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

// struct passed to the html template
type Article struct {
	Category        string // article category identifier
	Root            string // filepath to the article root location
	Id              string // human readable identifier
	Url             string // article path
	ReadmePath      string // filepath to README.md with article contents
	Summary         string // plain text summary
	WordCount       int
	ReadTimeMinutes int

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

func NewArticle(f string, c string, urlPrefix string, tags []string) (*Article, error) {
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

	const previewThreshold int = 256
	threshold := min(len(summary), previewThreshold)
	if threshold != len(summary) {
		summary = append(summary[:threshold-3], []byte("...")...)
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

	tagSet := make(map[string]bool)
	for _, v := range article.Tags {
		tagSet[v] = true
	}

	for _, v := range tags {
		if tagSet[v] {
			continue
		}
		article.Tags = append(article.Tags, v)
	}

	const readingWpm float64 = 200
	readTime, err := time.ParseDuration(
		strconv.Itoa(
			int(math.Ceil(float64(wordCount)/readingWpm))) + "m")

	if err != nil {
		readTime, _ = time.ParseDuration("2m")
	}

	article.ReadTimeMinutes = int(math.Ceil(readTime.Minutes()))

	return article, nil
}

func sanitizeUrl(s string) string {
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	return s + "/"
}
