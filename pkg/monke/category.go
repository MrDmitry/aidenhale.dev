package monke

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

type Category struct {
	Root            string              // filepath to the category root location
	Id              string              // human readable identifier
	Url             string              // article path
	ReadmePath      string              // filepath to README.md with category contents
	Articles        map[string]*Article // map of article by their id
	articlesCreated []*Article          // list of articles ordered by their created date, descending

	CategoryMeta
}

type CategoryMeta struct {
	Title string
}

func (c Category) GetArticlesByTime(limit int, offset int) []*Article {
	limit = min(limit, len(c.articlesCreated))
	in := c.articlesCreated[offset:limit]
	if len(in) == 0 {
		return nil
	}

	return in
}

func NewCategory(f string) (*Category, error) {
	dir, err := os.Open(f)
	if err != nil {
		log.Fatalf("Failed to open %s: %+v", f, err)
		return nil, err
	}
	defer dir.Close()

	contents, err := dir.ReadDir(0)
	if err != nil {
		log.Fatalf("Failed to read %s: %+v", f, err)
		return nil, err
	}

	var meta CategoryMeta

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

	readmePath := path.Join(f, "README.md")
	_, err = os.Open(readmePath)
	if err != nil {
		log.Warnf("Failed to open %s: %+v", readmePath, err)
		readmePath = ""
	}

	category := new(Category)
	category.Root = f
	category.Id = filepath.Base(category.Root)
	category.Url = fmt.Sprintf("/blog/%s/", category.Id)
	category.ReadmePath = readmePath
	category.Articles = make(map[string]*Article)
	category.CategoryMeta = meta

	for _, item := range contents {
		if !item.IsDir() {
			continue
		}
		cName := item.Name()
		article, err := NewArticle(filepath.Join(f, cName), category)
		if err != nil {
			continue
		}
		category.Articles[cName] = article
		category.articlesCreated = append(category.articlesCreated, article)
	}

	return category, nil
}
