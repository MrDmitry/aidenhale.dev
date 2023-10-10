package monke

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

type Monkebase struct {
	Categories map[string]*Category
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

	for _, item := range contents {
		if !item.IsDir() {
			continue
		}

		category, err := NewCategory(filepath.Join(r, item.Name()))
		if err != nil {
			log.Warnf("Failed to load category from %s: %+v", item.Name(), err)
			return err
		}

		Db.Categories[item.Name()] = category
	}

	return nil
}
