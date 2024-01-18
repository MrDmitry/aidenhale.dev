package main

import (
	"flag"
	"fmt"
	"mrdmitry/blog/pkg/monke"
	"os"
)

func main() {
	// find git directory for `git` calls
	workdir, err := os.Getwd()
	if err != nil {
		panic("failed to detect current work directory, aborting")
	}
	gitdir := workdir + "/.git"
	flag.Func("git-dir", "path to the website's .git directory (default `$PWD/.git`)", func(value string) error {
		gitdir, err = monke.DirectoryValidator(value)
		return err
	})

	protocol := "https"
	flag.Func("protocol", "webserver protocol: [http, https] (default \"https\")", func(value string) error {
		protocol, err = monke.ProtocolValidator(value)
		return err
	})

	hostname := flag.String("hostname", "aidenhale.dev", "hostname")

	flag.Parse()

	err = monke.InitDb("./web/data")

	// validate and set gitdir
	err = monke.GitDirectoryValidator(gitdir)
	if err != nil {
		panic(err)
	}
	monke.Gitdir = gitdir

	urlPrefix := monke.SanitizeUrl(fmt.Sprintf("%s://%s/", protocol, *hostname))

	data, err := monke.SitemapXml(urlPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println(string(data))
}
