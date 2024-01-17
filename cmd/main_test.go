package main

import (
	"testing"
)

func TestTrailingSlashNeeded(t *testing.T) {
	// everything under dir routers is considered to be a file
	dirRouters := []Router{
		{
			"/dir/": "dummy",
		},
	}

	// file routers contain verbatim path
	fileRouters := []Router{
		{
			"/robots.txt": "dummy",
		},
	}

	type testEntry struct {
		path    string
		want    bool
		message string
	}

	testData := []testEntry{
		// Dir-based match
		{path: "/dir", want: true, message: "router directory match"},
		{path: "/dir/", want: true, message: "router directory match"},
		{path: "/dir/filename", want: false, message: "router match of directory contents"},
		{path: "/dir/dirname/", want: false, message: "router match of directory contents"},
		{path: "/dirX", want: true, message: "typo in directory name"},
		{path: "/dirX/filename", want: true, message: "typo in directory name"},
		{path: "/root/dir/filename", want: true, message: "in-path match"},
		{path: "/root/dir/dirname/", want: true, message: "in-path match"},

		// File-based match
		{path: "/robots.txt", want: false, message: "router filepath match"},
		{path: "/robots.txt/", want: false, message: "router filepath match"},
		{path: "/robots.txt/etc", want: true, message: "ambiguous directory name"},
		{path: "/root/robots.txt", want: true, message: "subdirectory filepath match"},
		{path: "/root/robots.txt/etc", want: true, message: "ambiguous subdirectory name"},
		{path: "/robots.txt.bad", want: true, message: "filename with a suffix"},

		// Hardcoded router bit
		{path: "/assets", want: true, message: "intended usage"},
		{path: "/assets/", want: true, message: "intended usage"},
		{path: "/assets/filename", want: false, message: "intended usage"},
		{path: "/assets/dirname/", want: false, message: "assets hold only files"},
		{path: "/assets/assets/", want: false, message: "assets hold only files"},
		{path: "/assets/assets/filename", want: false, message: "stacked subdirectory assets contents"},
		{path: "/assets/assets/dirname/", want: false, message: "stacked subdirectory assets contents"},
		{path: "/root/assets/filename", want: false, message: "subdirectory assets match"},
		{path: "/root/assets/dirname/", want: false, message: "assets hold only files"},
		{path: "/root/assets/dirname/filename", want: true, message: "in-path assets"},

		// Misfits
		{path: "", want: true, message: "empty input"},
		{path: "a", want: true, message: "relative path"},
		{path: "/", want: true, message: "root"},
	}

	for _, entry := range testData {
		got := TrailingSlashNeeded(entry.path, dirRouters, fileRouters)
		if got != entry.want {
			t.Fatalf("%s: %s; expected %v, got %v", entry.path, entry.message, entry.want, got)
		}
	}
}
