package monke

import (
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

var Gitdir string = "./.git"

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

func GitRevision() (string, error) {
	cmd := exec.Command("git", "--git-dir", Gitdir, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		log.Warnf("failed to detect git revision: %+v", err)
		return strconv.FormatInt(time.Now().Unix(), 10), err
	}
	return strings.Trim(string(out), "\n"), nil
}

func GitLastLogTimeISO(path string) (time.Time, error) {
	cmd := exec.Command("git", "--git-dir", Gitdir, "log", "-1", "--pretty=%aI", "--", path)
	out, err := cmd.Output()
	if err != nil {
		log.Warnf("failed to detect git log for %s: %+v", path, err)
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, strings.Trim(string(out), "\n"))
}
