package monke

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

type Monkebase struct {
	Categories map[string]*Category
	Articles   ArticleLookup
}

var Db *Monkebase

func InitDb(r string) error {
	if Db != nil {
		return errors.New("Db already initialized")
	}

	Db = new(Monkebase)
	Db.Categories = make(map[string]*Category)

	root, err := os.Open(r)
	if err != nil {
		log.Fatalf("Failed to load database from %s: %+v", r, err)
		return err
	}
	defer root.Close()

	contents, err := root.ReadDir(0)
	if err != nil {
		log.Fatalf("Failed to load database from %s: %+v", r, err)
		return err
	}

	articles := make([]*Article, 0, 8)

	for _, item := range contents {
		if !item.IsDir() {
			continue
		}

		categoryPath := filepath.Join(r, item.Name())
		category, err := NewCategory(categoryPath)
		if err != nil {
			log.Warnf("Failed to load category from %s: %+v", item.Name(), err)
			continue
		}

		Db.Categories[item.Name()] = category

		dir, err := os.Open(categoryPath)
		if err != nil {
			log.Fatalf("Failed to open %s: %+v", categoryPath, err)
			return err
		}
		defer dir.Close()

		dirContents, err := dir.ReadDir(0)
		if err != nil {
			log.Fatalf("Failed to read contents of %s: %+v", categoryPath, err)
			return err
		}

		for _, item := range dirContents {
			if !item.IsDir() {
				continue
			}
			articleId := item.Name()
			articlePath := filepath.Join(categoryPath, articleId)
			article, err := NewArticle(articlePath, category.Id, category.Url, category.Tags)
			if err != nil {
				log.Warnf("Failed to process article %s", articlePath)
				continue
			}
			articles = append(articles, article)
		}
	}

	Db.Articles.Init(articles)

	if len(Db.Categories) == 0 {
		return errors.New("failed to load categories")
	}

	return nil
}
