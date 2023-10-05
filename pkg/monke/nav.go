package monke

import (
    "log"
    "os"
    "fmt"
    "path"
    "path/filepath"

    toml "github.com/pelletier/go-toml"
)

type NavItem struct {
    Url string
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

    Nav.Items = []NavItem {}

    for _, item := range files {
        file, err := os.Open(item)
        if err != nil {
            log.Fatalf("could not access %s: %+v", item, err)
        }

        info, err := file.Stat()
        if err != nil {
            log.Fatalf("could not stat %s: %+v", item, err)
        }

        if !info.IsDir() {
            continue
        }

        var meta NavMeta

        monkeMeta := path.Join(item, "monke.toml")
        tree, err := toml.LoadFile(monkeMeta)

        if err != nil {
            log.Fatalf("could not process %s: %+v", monkeMeta, err)
        }

        err = tree.Unmarshal(&meta)

        if err != nil {
            log.Fatalf("could not unwrap %s: %+v", monkeMeta, err)
        }

        Nav.Items = append(Nav.Items, NavItem {
            Url: fmt.Sprintf("/blog/%s", info.Name()),
            Name: meta.Title,
        })
    }
}
