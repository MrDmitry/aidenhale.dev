package monke

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/gommon/log"
	toml "github.com/pelletier/go-toml"
)

type NavItem struct {
	Url  string
	Name string
}

type NavData struct {
	Items []NavItem
}

type NavMeta struct {
	Title string
}

var Nav NavData

func NavInit() {
	globPattern := "./web/data/*"
	files, err := filepath.Glob(globPattern)

	if err != nil {
		log.Fatalf("could not find navigation data under %s: %+v", globPattern, err)
	}

	Nav.Items = []NavItem{}

	for _, item := range files {
		file, err := os.Open(item)
		if err != nil {
			log.Warnf("could not access %s, skipping: %+v", item, err)
			continue
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			log.Warnf("could not stat %s, skipping: %+v", item, err)
			continue
		}

		if !info.IsDir() {
			continue
		}

		var meta NavMeta

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

		Nav.Items = append(Nav.Items, NavItem{
			Url:  fmt.Sprintf("/blog/%s/", info.Name()),
			Name: meta.Title,
		})
	}
}
