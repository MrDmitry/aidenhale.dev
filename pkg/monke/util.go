package monke

import (
	"net/url"
	"path/filepath"
	"strings"
)

func sanitizePath(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] != '/' {
		s = "/" + s
	}
	trailingSlash := s[len(s)-1] == '/'
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	s = strings.TrimRight(s, "/")

	if trailingSlash {
		return s + "/"
	}
	return s
}

func SanitizeUrl(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	u.Path = sanitizePath(u.Path)
	return u.String()
}
