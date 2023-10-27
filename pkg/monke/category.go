package monke

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

type Category struct {
	Root       string // filepath to the category root location
	Id         string // human readable identifier
	Url        string // article path
	ReadmePath string // filepath to README.md with category contents

	CategoryData
}

type categoryMeta struct {
	Version int
}

type CategoryData struct {
	Title string
	Tags  []string
}

type CategoryToml struct {
	Meta categoryMeta
	Data CategoryData
}

func NewCategory(f string) (*Category, error) {
	dir, err := os.Open(f)
	if err != nil {
		log.Fatalf("failed to open %s: %+v", f, err)
		return nil, err
	}
	defer dir.Close()

	var meta CategoryToml

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
		return nil, errors.New("incompatible version")
	}

	readmePath := path.Join(f, "README.md")
	_, err = os.Open(readmePath)
	if err != nil {
		log.Warnf("failed to open %s: %+v", readmePath, err)
		readmePath = ""
	}

	category := new(Category)
	category.Root = f
	category.Id = filepath.Base(category.Root)
	category.Url = fmt.Sprintf("/blog/%s/", category.Id)
	category.ReadmePath = readmePath
	category.CategoryData = meta.Data

	return category, nil
}
