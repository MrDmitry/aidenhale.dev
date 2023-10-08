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

type Article struct {
	Id      string
	Url     string
	Summary string

	ArticleMeta
}

type ArticleMeta struct {
	Title   string
	Created time.Time
}

func GetArticles(topic string) ([]Article, error) {
	globPattern := fmt.Sprintf("./web/data/%s/*", topic)
	entries, err := filepath.Glob(globPattern)

	if err != nil {
		log.Warnf("failed to find anything under %s: %+v", globPattern, err)
		return nil, err
	}

	var results []Article
	for _, item := range entries {
		entry, err := os.Open(item)
		if err != nil {
			log.Debugf("failed to open %s, skipping: %+v", entry, err)
			continue
		}
		defer entry.Close()

		info, err := entry.Stat()
		if err != nil {
			log.Debugf("failed to stat %s, skipping: %+v", entry, err)
			continue
		}
		if !info.IsDir() {
			continue
		}

		var meta ArticleMeta

		monkeMeta := path.Join(item, "monke.toml")
		tree, err := toml.LoadFile(monkeMeta)

		if err != nil {
			log.Warnf("could not process %s, skipping: %+v", monkeMeta, err)
			continue
		}

		err = tree.Unmarshal(&meta)

		if err != nil {
			log.Warnf("could not unwrap %s, skipping: %+v", monkeMeta, err)
			continue
		}

		readme := path.Join(item, "README.md")
		summary, err := RenderMarkdownToTextPreview(readme, 200)

		if err != nil {
			log.Warnf("could not process %s, skipping: %+v", readme, err)
			continue
		}

		results = append(results, Article{
			Id:          info.Name(),
			Url:         fmt.Sprintf("/blog/%s/%s/", topic, info.Name()),
			Summary:     string(summary),
			ArticleMeta: meta,
		})
	}

	return results, nil
}
