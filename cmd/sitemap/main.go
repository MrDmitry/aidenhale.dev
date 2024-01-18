package main

import (
	"flag"
	"fmt"
	"mrdmitry/blog/pkg/monke"
)

func main() {
	gitdir := monke.Config.Gitdir
	flag.Func("git-dir", fmt.Sprintf("path to the website's .git directory (default `%s`)", gitdir), func(value string) error {
		v, err := monke.DirectoryValidator(value)
		gitdir = v
		return err
	})

	protocol := "https"
	flag.Func("protocol", "webserver protocol: [http, https] (default \"https\")", func(value string) error {
		v, err := monke.ProtocolValidator(value)
		protocol = v
		return err
	})

	hostname := flag.String("hostname", "aidenhale.dev", "hostname")

	flag.Parse()

	monke.Config.Gitdir = gitdir
	if err := monke.Config.Validate(); err != nil {
		panic(err)
	}

	if err := monke.InitDb("./web/data"); err != nil {
		panic(err)
	}

	urlPrefix := monke.SanitizeUrl(fmt.Sprintf("%s://%s/", protocol, *hostname))

	data, err := monke.SitemapXml(urlPrefix)
	if err != nil {
		panic(err)
	}

	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println(string(data))
}
