package monke

import (
	"net/url"
	"path/filepath"
)

func SanitizePath(s string) string {
	if len(s) == 1 {
		return s
	}
	trailingSlash := s[len(s)-1] == '/'
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	switch trailingSlash {
	case true:
		return s + "/"
	default:
		return s
	}
}

func SanitizeUrl(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	u.Path = SanitizePath(u.Path)
	return u.String()
}
